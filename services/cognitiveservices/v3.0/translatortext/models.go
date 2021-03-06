package translatortext

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
)

// The package's fully qualified name.
const fqdn = "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/translatortext"

// BreakSentenceResultItem ...
type BreakSentenceResultItem struct {
	SentLen *[]int32 `json:"sentLen,omitempty"`
}

// BreakSentenceTextInput text needed for break sentence request
type BreakSentenceTextInput struct {
	Text *string `json:"text,omitempty"`
}

// DetectResultItem ...
type DetectResultItem struct {
	Text *string `json:"text,omitempty"`
}

// DetectTextInput text needed for detect request
type DetectTextInput struct {
	Text *string `json:"text,omitempty"`
}

// DictionaryExampleResultItem ...
type DictionaryExampleResultItem struct {
	NormalizedSource *string                                    `json:"normalizedSource,omitempty"`
	NormalizedTarget *string                                    `json:"normalizedTarget,omitempty"`
	Examples         *[]DictionaryExampleResultItemExamplesItem `json:"examples,omitempty"`
}

// DictionaryExampleResultItemExamplesItem ...
type DictionaryExampleResultItemExamplesItem struct {
	SourcePrefix *string `json:"sourcePrefix,omitempty"`
	SourceTerm   *string `json:"sourceTerm,omitempty"`
	SourceSuffix *string `json:"sourceSuffix,omitempty"`
	TargetPrefix *string `json:"targetPrefix,omitempty"`
	TargetTerm   *string `json:"targetTerm,omitempty"`
	TargetSuffix *string `json:"targetSuffix,omitempty"`
}

// DictionaryExampleTextInput text needed for a dictionary example request
type DictionaryExampleTextInput struct {
	Text        *string `json:"text,omitempty"`
	Translation *string `json:"translation,omitempty"`
}

// DictionaryLookupResultItem ...
type DictionaryLookupResultItem struct {
	NormalizedSource *string                                       `json:"normalizedSource,omitempty"`
	DisplaySource    *string                                       `json:"displaySource,omitempty"`
	Translations     *[]DictionaryLookupResultItemTranslationsItem `json:"translations,omitempty"`
}

// DictionaryLookupResultItemTranslationsItem ...
type DictionaryLookupResultItemTranslationsItem struct {
	NormalizedTarget *string                                                           `json:"normalizedTarget,omitempty"`
	DisplayTarget    *string                                                           `json:"displayTarget,omitempty"`
	PosTag           *string                                                           `json:"posTag,omitempty"`
	Confidence       *float64                                                          `json:"confidence,omitempty"`
	PrefixWord       *string                                                           `json:"prefixWord,omitempty"`
	BackTranslations *[]DictionaryLookupResultItemTranslationsItemBackTranslationsItem `json:"backTranslations,omitempty"`
}

// DictionaryLookupResultItemTranslationsItemBackTranslationsItem ...
type DictionaryLookupResultItemTranslationsItemBackTranslationsItem struct {
	NormalizedText *string `json:"normalizedText,omitempty"`
	DisplayText    *string `json:"displayText,omitempty"`
	NumExamples    *int32  `json:"numExamples,omitempty"`
	FrequencyCount *int32  `json:"frequencyCount,omitempty"`
}

// DictionaryLookupTextInput text needed for a dictionary lookup request
type DictionaryLookupTextInput struct {
	Text *string `json:"text,omitempty"`
}

// ErrorMessage ...
type ErrorMessage struct {
	Error *ErrorMessageError `json:"error,omitempty"`
}

// ErrorMessageError ...
type ErrorMessageError struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

// LanguagesResult example of a successful languages request
type LanguagesResult struct {
	autorest.Response `json:"-"`
	Translation       *LanguagesResultTranslation     `json:"translation,omitempty"`
	Transliteration   *LanguagesResultTransliteration `json:"transliteration,omitempty"`
	Dictionary        *LanguagesResultDictionary      `json:"dictionary,omitempty"`
}

// LanguagesResultDictionary ...
type LanguagesResultDictionary struct {
	LanguageCode *LanguagesResultDictionaryLanguageCode `json:"languageCode,omitempty"`
}

// LanguagesResultDictionaryLanguageCode ...
type LanguagesResultDictionaryLanguageCode struct {
	Name         *string                                                  `json:"name,omitempty"`
	NativeName   *string                                                  `json:"nativeName,omitempty"`
	Dir          *string                                                  `json:"dir,omitempty"`
	Translations *[]LanguagesResultDictionaryLanguageCodeTranslationsItem `json:"translations,omitempty"`
}

// LanguagesResultDictionaryLanguageCodeTranslationsItem ...
type LanguagesResultDictionaryLanguageCodeTranslationsItem struct {
	Name       *string `json:"name,omitempty"`
	NativeName *string `json:"nativeName,omitempty"`
	Dir        *string `json:"dir,omitempty"`
	Code       *string `json:"code,omitempty"`
}

// LanguagesResultTranslation ...
type LanguagesResultTranslation struct {
	LanguageCode *LanguagesResultTranslationLanguageCode `json:"languageCode,omitempty"`
}

// LanguagesResultTranslationLanguageCode ...
type LanguagesResultTranslationLanguageCode struct {
	Name       *string `json:"name,omitempty"`
	NativeName *string `json:"nativeName,omitempty"`
	Dir        *string `json:"dir,omitempty"`
}

// LanguagesResultTransliteration ...
type LanguagesResultTransliteration struct {
	LanguageCode *LanguagesResultTransliterationLanguageCode `json:"languageCode,omitempty"`
}

// LanguagesResultTransliterationLanguageCode ...
type LanguagesResultTransliterationLanguageCode struct {
	Name       *string                                                  `json:"name,omitempty"`
	NativeName *string                                                  `json:"nativeName,omitempty"`
	Scripts    *[]LanguagesResultTransliterationLanguageCodeScriptsItem `json:"scripts,omitempty"`
}

// LanguagesResultTransliterationLanguageCodeScriptsItem ...
type LanguagesResultTransliterationLanguageCodeScriptsItem struct {
	Code       *string                                                               `json:"code,omitempty"`
	Name       *string                                                               `json:"name,omitempty"`
	NativeName *string                                                               `json:"nativeName,omitempty"`
	Dir        *string                                                               `json:"dir,omitempty"`
	ToScripts  *[]LanguagesResultTransliterationLanguageCodeScriptsItemToScriptsItem `json:"toScripts,omitempty"`
}

// LanguagesResultTransliterationLanguageCodeScriptsItemToScriptsItem ...
type LanguagesResultTransliterationLanguageCodeScriptsItemToScriptsItem struct {
	Code       *string `json:"code,omitempty"`
	Name       *string `json:"name,omitempty"`
	NativeName *string `json:"nativeName,omitempty"`
	Dir        *string `json:"dir,omitempty"`
}

// ListBreakSentenceResultItem ...
type ListBreakSentenceResultItem struct {
	autorest.Response `json:"-"`
	Value             *[]BreakSentenceResultItem `json:"value,omitempty"`
}

// ListDetectResultItem ...
type ListDetectResultItem struct {
	autorest.Response `json:"-"`
	Value             *[]DetectResultItem `json:"value,omitempty"`
}

// ListDictionaryExampleResultItem ...
type ListDictionaryExampleResultItem struct {
	autorest.Response `json:"-"`
	Value             *[]DictionaryExampleResultItem `json:"value,omitempty"`
}

// ListDictionaryLookupResultItem ...
type ListDictionaryLookupResultItem struct {
	autorest.Response `json:"-"`
	Value             *[]DictionaryLookupResultItem `json:"value,omitempty"`
}

// ListTranslateResultAllItem ...
type ListTranslateResultAllItem struct {
	autorest.Response `json:"-"`
	Value             *[]TranslateResultAllItem `json:"value,omitempty"`
}

// ListTransliterateResultItem ...
type ListTransliterateResultItem struct {
	autorest.Response `json:"-"`
	Value             *[]TransliterateResultItem `json:"value,omitempty"`
}

// TranslateResultAllItem ...
type TranslateResultAllItem struct {
	DetectedLanguage *TranslateResultAllItemDetectedLanguage   `json:"detectedLanguage,omitempty"`
	Translations     *[]TranslateResultAllItemTranslationsItem `json:"translations,omitempty"`
}

// TranslateResultAllItemDetectedLanguage ...
type TranslateResultAllItemDetectedLanguage struct {
	Language *string `json:"language,omitempty"`
	Score    *int32  `json:"score,omitempty"`
}

// TranslateResultAllItemTranslationsItem ...
type TranslateResultAllItemTranslationsItem struct {
	Text            *string                                                `json:"text,omitempty"`
	Transliteration *TranslateResultAllItemTranslationsItemTransliteration `json:"transliteration,omitempty"`
	To              *string                                                `json:"to,omitempty"`
	Alignment       *TranslateResultAllItemTranslationsItemAlignment       `json:"alignment,omitempty"`
	SentLen         *TranslateResultAllItemTranslationsItemSentLen         `json:"sentLen,omitempty"`
}

// TranslateResultAllItemTranslationsItemAlignment ...
type TranslateResultAllItemTranslationsItemAlignment struct {
	Proj *string `json:"proj,omitempty"`
}

// TranslateResultAllItemTranslationsItemSentLen ...
type TranslateResultAllItemTranslationsItemSentLen struct {
	SrcSentLen   *[]TranslateResultAllItemTranslationsItemSentLenSrcSentLenItem   `json:"srcSentLen,omitempty"`
	TransSentLen *[]TranslateResultAllItemTranslationsItemSentLenTransSentLenItem `json:"transSentLen,omitempty"`
}

// TranslateResultAllItemTranslationsItemSentLenSrcSentLenItem ...
type TranslateResultAllItemTranslationsItemSentLenSrcSentLenItem struct {
	Integer *int32 `json:"integer,omitempty"`
}

// TranslateResultAllItemTranslationsItemSentLenTransSentLenItem ...
type TranslateResultAllItemTranslationsItemSentLenTransSentLenItem struct {
	Integer *int32 `json:"integer,omitempty"`
}

// TranslateResultAllItemTranslationsItemTransliteration ...
type TranslateResultAllItemTranslationsItemTransliteration struct {
	Text   *string `json:"text,omitempty"`
	Script *string `json:"script,omitempty"`
}

// TranslateResultItem ...
type TranslateResultItem struct {
	Translation *[]TranslateResultItemTranslationItem `json:"translation,omitempty"`
}

// TranslateResultItemTranslationItem ...
type TranslateResultItemTranslationItem struct {
	Text *string `json:"text,omitempty"`
	To   *string `json:"to,omitempty"`
}

// TranslateTextInput text needed for a translate request
type TranslateTextInput struct {
	Text *string `json:"text,omitempty"`
}

// TransliterateResultItem ...
type TransliterateResultItem struct {
	Text   *string `json:"text,omitempty"`
	Script *string `json:"script,omitempty"`
}

// TransliterateTextInput text needed for a transliterate request
type TransliterateTextInput struct {
	Text *string `json:"text,omitempty"`
}
