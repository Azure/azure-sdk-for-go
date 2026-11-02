// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// cSpell:words usize

pub const DOCUMENT_BYTES: usize = 1024;
pub const SEED_COUNT: usize = 1024;

pub fn validate_owner(owner: &str) -> Result<(), String> {
    if !(8..=64).contains(&owner.len()) {
        return Err("owner must contain 8 through 64 characters".into());
    }
    if !owner
        .bytes()
        .all(|value| value.is_ascii_lowercase() || value.is_ascii_digit() || value == b'-')
    {
        return Err("owner may contain only lowercase ASCII letters, digits, and hyphens".into());
    }
    Ok(())
}

pub fn item_id(owner: &str, index: usize) -> String {
    format!("{owner}-{index:04}")
}

pub fn key_index(sequence: u64) -> usize {
    sequence as usize % SEED_COUNT
}

pub fn make_document(owner: &str, index: usize) -> Result<Vec<u8>, String> {
    validate_owner(owner)?;
    let id = item_id(owner, index);
    let prefix = format!("{{\"id\":\"{id}\",\"pk\":\"{id}\",\"owner\":\"{owner}\",\"padding\":\"");
    let suffix = "\"}";
    let padding = expected_padding_length(owner, index)?;
    let mut body = Vec::with_capacity(DOCUMENT_BYTES);
    body.extend_from_slice(prefix.as_bytes());
    body.extend(std::iter::repeat_n(b'x', padding));
    body.extend_from_slice(suffix.as_bytes());
    if body.len() != DOCUMENT_BYTES {
        return Err("document length invariant failed".into());
    }
    Ok(body)
}

pub fn validate_document(body: &[u8], owner: &str, index: usize) -> Result<(), String> {
    let value: serde_json::Value =
        serde_json::from_slice(body).map_err(|error| format!("decode document: {error}"))?;
    let id = item_id(owner, index);
    if value["id"] != id || value["pk"] != id || value["owner"] != owner {
        return Err("document id, pk, or owner does not match the campaign".into());
    }
    let padding = value["padding"]
        .as_str()
        .ok_or("document padding is missing or is not a string")?;
    if padding.len() != expected_padding_length(owner, index)?
        || !padding.bytes().all(|value| value == b'x')
    {
        return Err("document padding does not match the exact-size campaign payload".into());
    }
    Ok(())
}

fn expected_padding_length(owner: &str, index: usize) -> Result<usize, String> {
    const EMPTY_DOCUMENT_BYTES: usize = r#"{"id":"","pk":"","owner":"","padding":""}"#.len();
    let id = item_id(owner, index);
    DOCUMENT_BYTES
        .checked_sub(EMPTY_DOCUMENT_BYTES + 2 * id.len() + owner.len())
        .ok_or_else(|| "owner leaves no room for document padding".to_string())
}

pub fn owned_by(body: &[u8], owner: &str) -> bool {
    let Ok(text) = std::str::from_utf8(body) else {
        return false;
    };
    text.contains(&format!("\"owner\":\"{owner}\""))
}

pub fn reject_secret_inputs(arguments: &[String]) -> Result<(), String> {
    for argument in arguments {
        let lower = argument.to_ascii_lowercase();
        if matches!(lower.as_str(), "--key" | "--account-key")
            || lower.starts_with("--key=")
            || lower.starts_with("--account-key=")
        {
            return Err(format!(
                "secret command-line argument {argument:?} is forbidden"
            ));
        }
    }
    for name in [
        "AZURE_COSMOS_KEY",
        "COSMOS_KEY",
        "COSMOS_ACCOUNT_KEY",
        "ACCOUNT_KEY",
    ] {
        if std::env::var_os(name).is_some() {
            return Err(format!(
                "secret environment variable {name} is forbidden; use --key-fd"
            ));
        }
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn document_has_exact_size() {
        let body = make_document("campaign-owner", 17).unwrap();
        assert_eq!(body.len(), DOCUMENT_BYTES);
        validate_document(&body, "campaign-owner", 17).unwrap();
    }

    #[test]
    fn validation_allows_cosmos_system_properties() {
        let mut body = make_document("campaign-owner", 17).unwrap();
        body.pop();
        body.extend_from_slice(br#","_rid":"example","_etag":"example"}"#);
        assert!(body.len() > DOCUMENT_BYTES);
        validate_document(&body, "campaign-owner", 17).unwrap();
    }

    #[test]
    fn validation_rejects_changed_campaign_payload() {
        let body = make_document("campaign-owner", 17).unwrap();
        let mut value: serde_json::Value = serde_json::from_slice(&body).unwrap();
        value["padding"] =
            serde_json::Value::String("y".repeat(value["padding"].as_str().unwrap().len()));
        let changed = serde_json::to_vec(&value).unwrap();
        assert!(validate_document(&changed, "campaign-owner", 17).is_err());
    }

    #[test]
    fn key_sequence_is_round_robin() {
        assert_eq!(key_index(0), 0);
        assert_eq!(key_index(1023), 1023);
        assert_eq!(key_index(1024), 0);
    }

    #[test]
    fn owner_guard_rejects_other_owner() {
        let body = make_document("campaign-owner", 1).unwrap();
        assert!(!owned_by(&body, "other-owner"));
    }

    #[test]
    fn command_line_secret_is_rejected() {
        assert!(reject_secret_inputs(&["run".into(), "--key=secret".into()]).is_err());
    }
}
