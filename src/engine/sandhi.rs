use fancy_regex::Regex;
use lazy_static::lazy_static;
use crate::engine::declension::apply_natva;

struct UpasargaRule {
    pattern: Regex,
    canonical: &'static str,
    prepend: &'static str,
}

lazy_static! {
    static ref UPASARGA_PATTERNS: Vec<UpasargaRule> = vec![
        UpasargaRule { pattern: Regex::new(r"^samudA").unwrap(), canonical: "sam + ud + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^sa[mMnYNR]u[dtcjNYRnl]A").unwrap(), canonical: "sam + ud + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^sa[mMnYNR]u[dtcjNYRnl]").unwrap(), canonical: "sam + ud", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^vyA").unwrap(), canonical: "vi + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^pratyA").unwrap(), canonical: "prati + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^vy").unwrap(), canonical: "vi", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^vi").unwrap(), canonical: "vi", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^praty").unwrap(), canonical: "prati", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^prati").unwrap(), canonical: "prati", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^sa[mMnYNR]").unwrap(), canonical: "sam", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^sam").unwrap(), canonical: "sam", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^pra").unwrap(), canonical: "pra", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^A").unwrap(), canonical: "AN", prepend: "" },
    ];
}

pub fn get_upasarga_splits(word: &str) -> Vec<(String, String)> {
    let mut splits = Vec::new();
    for rule in UPASARGA_PATTERNS.iter() {
        if let Ok(Some(mat)) = rule.pattern.find(word) {
            let end_idx = mat.end();
            let stripped = format!("{}{}", rule.prepend, &word[end_idx..]);
            if stripped.chars().count() >= 2 {
                splits.push((rule.canonical.to_string(), stripped));
            }
        }
    }
    splits
}

pub fn apply_upasarga_sandhi(upasarga_str: &str, form: &str) -> String {
    if upasarga_str.is_empty() { return form.to_string(); }
    let prefixes: Vec<&str> = upasarga_str.split('+').map(|s| s.trim()).collect();
    let mut result = form.to_string();
    
    for mut p in prefixes.into_iter().rev() {
        if p == "AN" { p = "A"; }
        
        let p_len = p.len();
        let p_last = p.chars().last().unwrap_or(' ');
        let r_first = result.chars().next().unwrap_or(' ');
        let is_vowel = |c: char| "aAiIuUfFeEoO".contains(c);
        
        if p_last == 'i' && is_vowel(r_first) {
            result = format!("{}y{}", &p[..p_len-1], result);
        } else if p_last == 'u' && is_vowel(r_first) {
            result = format!("{}v{}", &p[..p_len-1], result);
        } else if p_last == 'a' || p_last == 'A' {
            if r_first == 'a' || r_first == 'A' { result = format!("{}A{}", &p[..p_len-1], &result[1..]); }
            else if r_first == 'i' || r_first == 'I' { result = format!("{}e{}", &p[..p_len-1], &result[1..]); }
            else if r_first == 'u' || r_first == 'U' { result = format!("{}o{}", &p[..p_len-1], &result[1..]); }
            else if r_first == 'f' || r_first == 'F' { result = format!("{}ar{}", &p[..p_len-1], &result[1..]); }
            else { result = format!("{}{}", p, result); }
        } else if p_last == 'm' {
            if !is_vowel(r_first) { result = format!("{}M{}", &p[..p_len-1], result); }
            else { result = format!("{}{}", p, result); }
        } else if p_last == 'd' {
            if "kKqQpPzSs".contains(r_first) { result = format!("{}t{}", &p[..p_len-1], result); }
            else if "cC".contains(r_first) { result = format!("{}c{}", &p[..p_len-1], result); }
            else if "jJ".contains(r_first) { result = format!("{}j{}", &p[..p_len-1], result); }
            else { result = format!("{}{}", p, result); }
        } else {
            result = format!("{}{}", p, result);
        }
        
        result = apply_natva(&result);
    }
    result
}
