use fancy_regex::Regex;
use lazy_static::lazy_static;
use serde_json::{json, Value};

struct NamadhatuRule {
    pattern: Regex,
    pratyaya: &'static str,
    fallback_vowel: &'static str,
}

lazy_static! {
    static ref NAMADHATU_RULES: Vec<NamadhatuRule> = vec![
        NamadhatuRule { pattern: Regex::new(r"Iya(ti|taH|nti|si|TaH|Ta|mi|vaH|maH|te|yate|yante)$").unwrap(), pratyaya: "kyac", fallback_vowel: "a" },
        NamadhatuRule { pattern: Regex::new(r"Aya(te|yete|yante|se|yeTe|yaDve|ye|yAvahe|yAmahe)$").unwrap(), pratyaya: "kyaN", fallback_vowel: "a" },
        NamadhatuRule { pattern: Regex::new(r"kAmya(ti|taH|nti|si|TaH|Ta|mi|vaH|maH)$").unwrap(), pratyaya: "kAmyac", fallback_vowel: "" },
    ];
}

pub fn analyze_namadhatu(word: &str) -> Vec<Value> {
    let mut results = Vec::new();
    for rule in NAMADHATU_RULES.iter() {
        if rule.pattern.is_match(word).unwrap_or(false) {
            let stem = rule.pattern.replace(word, "").into_owned();
            
            let base_noun = if (rule.pratyaya == "kyac" || rule.pratyaya == "kyaN") && !rule.fallback_vowel.is_empty() {
                format!("{}{}", stem, rule.fallback_vowel)
            } else {
                stem
            };

            results.push(json!({
                "type": "namadhatu",
                "base_noun": base_noun,
                "pratyaya": rule.pratyaya,
                "meaning_hint": format!("Desires/behaves like {}", base_noun),
                "note": format!("Derived algorithmically using {} suffix", rule.pratyaya)
            }));
        }
    }
    results
}
