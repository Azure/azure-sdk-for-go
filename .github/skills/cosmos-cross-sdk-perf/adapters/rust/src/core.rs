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
    let padding = DOCUMENT_BYTES
        .checked_sub(prefix.len() + suffix.len())
        .ok_or_else(|| "owner leaves no room for document padding".to_string())?;
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
    if body.len() != DOCUMENT_BYTES {
        return Err(format!(
            "document length is {}, expected {DOCUMENT_BYTES}",
            body.len()
        ));
    }
    let text = std::str::from_utf8(body).map_err(|error| error.to_string())?;
    let id = item_id(owner, index);
    for expected in [
        format!("\"id\":\"{id}\""),
        format!("\"pk\":\"{id}\""),
        format!("\"owner\":\"{owner}\""),
    ] {
        if !text.contains(&expected) {
            return Err("document id, pk, or owner does not match the campaign".into());
        }
    }
    Ok(())
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
