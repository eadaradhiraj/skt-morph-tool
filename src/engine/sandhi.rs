use fancy_regex::Regex;
use lazy_static::lazy_static;

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
        UpasargaRule { pattern: Regex::new(r"^sa[mMnYNR]praty").unwrap(), canonical: "sam + prati", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^sa[mMnYNR]prati").unwrap(), canonical: "sam + prati", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^sa[mMnYNR]pra").unwrap(), canonical: "sam + pra", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^vyA").unwrap(), canonical: "vi + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^pratyA").unwrap(), canonical: "prati + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^paryA").unwrap(), canonical: "pari + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^udA").unwrap(), canonical: "ud + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^prA").unwrap(), canonical: "pra + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^apA").unwrap(), canonical: "apa + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^anvA").unwrap(), canonical: "anu + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^avA").unwrap(), canonical: "ava + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^nirA").unwrap(), canonical: "nir + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^durA").unwrap(), canonical: "dus + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^aDyA").unwrap(), canonical: "aDi + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^apyA").unwrap(), canonical: "api + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^atyA").unwrap(), canonical: "ati + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^aByA").unwrap(), canonical: "aBi + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^upA").unwrap(), canonical: "upa + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^svA").unwrap(), canonical: "su + AN", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^vy").unwrap(), canonical: "vi", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^vi").unwrap(), canonical: "vi", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^praty").unwrap(), canonical: "prati", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^prati").unwrap(), canonical: "prati", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^pary").unwrap(), canonical: "pari", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^pari").unwrap(), canonical: "pari", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^sa[mMnYNR]").unwrap(), canonical: "sam", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^sam").unwrap(), canonical: "sam", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^u[dt]y").unwrap(), canonical: "ud", prepend: "y" },
        UpasargaRule { pattern: Regex::new(r"^u[dtcjNYRnl]").unwrap(), canonical: "ud", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^pra").unwrap(), canonical: "pra", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^parA").unwrap(), canonical: "parA", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^apa").unwrap(), canonical: "apa", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^anv").unwrap(), canonical: "anu", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^anu").unwrap(), canonical: "anu", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^ava").unwrap(), canonical: "ava", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^ni[rsSzR]").unwrap(), canonical: "nir", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^du[rsSzR]").unwrap(), canonical: "dus", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^ny").unwrap(), canonical: "ni", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^ni").unwrap(), canonical: "ni", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^aDy").unwrap(), canonical: "aDi", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^aDi").unwrap(), canonical: "aDi", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^apy").unwrap(), canonical: "api", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^api").unwrap(), canonical: "api", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^aty").unwrap(), canonical: "ati", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^ati").unwrap(), canonical: "ati", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^sv").unwrap(), canonical: "su", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^su").unwrap(), canonical: "su", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^aBy").unwrap(), canonical: "aBi", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^aBi").unwrap(), canonical: "aBi", prepend: "" },
        UpasargaRule { pattern: Regex::new(r"^upa").unwrap(), canonical: "upa", prepend: "" },
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
