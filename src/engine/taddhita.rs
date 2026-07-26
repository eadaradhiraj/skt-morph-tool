use fancy_regex::Regex;
use lazy_static::lazy_static;
use serde_json::{json, Value};

fn reverse_vriddhi(word: &str) -> String {
    if word.is_empty() { return word.to_string(); }
    
    let mut chars: Vec<char> = word.chars().collect();
    
    // Check if the word starts with a vowel
    if chars[0] == 'A' { chars[0] = 'a'; return chars.into_iter().collect(); }
    if chars[0] == 'E' { chars[0] = 'i'; return chars.into_iter().collect(); }
    if chars[0] == 'O' { chars[0] = 'u'; return chars.into_iter().collect(); }

    // If starts with consonant, find the first vowel
    if let Ok(Some(mat)) = Regex::new(r"[aAiIuUfFeEoO]").unwrap().find(word) {
        let idx = mat.start();
        // Handle "ar" specifically if needed, but normally single char replace is enough
        let v = chars[idx];
        if v == 'A' { chars[idx] = 'a'; }
        else if v == 'E' { chars[idx] = 'i'; }
        else if v == 'O' { chars[idx] = 'u'; }
    }
    
    chars.into_iter().collect()
}

struct TaddhitaRule {
    pattern: Regex,
    pratyaya: &'static str,
    fallback_vowel: &'static str,
    apply_vriddhi: bool,
    meaning: &'static str,
}

lazy_static! {
    static ref TADDHITA_RULES: Vec<TaddhitaRule> = vec![
        TaddhitaRule { pattern: Regex::new(r"i$").unwrap(), pratyaya: "iY", fallback_vowel: "a", apply_vriddhi: true, meaning: "Descendant of " },
        TaddhitaRule { pattern: Regex::new(r"eya$").unwrap(), pratyaya: "Qak", fallback_vowel: "I", apply_vriddhi: true, meaning: "Descendant of " },
        TaddhitaRule { pattern: Regex::new(r"Aya[nR]a$").unwrap(), pratyaya: "PaY", fallback_vowel: "a", apply_vriddhi: true, meaning: "Descendant of " },
        TaddhitaRule { pattern: Regex::new(r"a$").unwrap(), pratyaya: "aR", fallback_vowel: "a", apply_vriddhi: true, meaning: "Descendant/Relation of " },
        
        TaddhitaRule { pattern: Regex::new(r"tva$").unwrap(), pratyaya: "tva", fallback_vowel: "", apply_vriddhi: false, meaning: "State/nature of " },
        TaddhitaRule { pattern: Regex::new(r"tA$").unwrap(), pratyaya: "tal", fallback_vowel: "", apply_vriddhi: false, meaning: "State/nature of " },
        TaddhitaRule { pattern: Regex::new(r"mat$").unwrap(), pratyaya: "matup", fallback_vowel: "", apply_vriddhi: false, meaning: "Possessing " },
        TaddhitaRule { pattern: Regex::new(r"vat$").unwrap(), pratyaya: "vatup", fallback_vowel: "", apply_vriddhi: false, meaning: "Possessing " },
        TaddhitaRule { pattern: Regex::new(r"in$").unwrap(), pratyaya: "ini", fallback_vowel: "", apply_vriddhi: false, meaning: "Possessing " },
    ];
}

pub fn analyze_taddhita(word: &str) -> Vec<Value> {
    let mut results = Vec::new();
    for rule in TADDHITA_RULES.iter() {
        if rule.pattern.is_match(word).unwrap_or(false) {
            let mut stem = rule.pattern.replace(word, "").into_owned();
            
            if rule.apply_vriddhi {
                stem = reverse_vriddhi(&stem);
            }
            
            let mut base_noun = format!("{}{}", stem, rule.fallback_vowel);
            
            if base_noun.ends_with("aa") {
                base_noun.pop();
            }

            results.push(json!({
                "type": "taddhita",
                "base_noun": base_noun,
                "pratyaya": rule.pratyaya,
                "meaning_hint": format!("{}{}", rule.meaning, base_noun)
            }));
        }
    }
    results
}
