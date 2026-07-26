use fancy_regex::Regex;
use lazy_static::lazy_static;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct GuessedStem {
    pub stem: String,
    pub gender: String,
    pub case: String,
    pub vacana: String,
}

struct StemRule {
    pattern: Regex,
    replacement: &'static str,
    gender: &'static str,
    case: &'static str,
    vacana: &'static str,
}

macro_rules! rule {
    ($pat:expr, $rep:expr, $gen:expr, $case:expr, $vac:expr) => {
        StemRule { pattern: Regex::new($pat).unwrap(), replacement: $rep, gender: $gen, case: $case, vacana: $vac }
    };
}

lazy_static! {
    static ref STEMMING_RULES: Vec<StemRule> = vec![
        rule!(r"aH$", "a", "masc", "prathama", "eka"), rule!(r"O$", "a", "masc", "prathama/dvitiya", "dvi"), rule!(r"AH$", "a", "masc", "prathama", "bahu"),
        rule!(r"am$", "a", "masc/neut", "dvitiya/prathama", "eka"), rule!(r"An$", "a", "masc", "dvitiya", "bahu"), rule!(r"e[nR]a$", "a", "masc/neut", "tritiya", "eka"),
        rule!(r"AByAm$", "a", "masc/neut", "tritiya/caturthi/panchami", "dvi"), rule!(r"EH$", "a", "masc/neut", "tritiya", "bahu"),
        rule!(r"Aya$", "a", "masc/neut", "caturthi", "eka"), rule!(r"eByaH$", "a", "masc/neut", "caturthi/panchami", "bahu"),
        rule!(r"At$", "a", "masc/neut", "panchami", "eka"), rule!(r"asya$", "a", "masc/neut", "sasthi", "eka"),
        rule!(r"ayoH$", "a", "masc/neut", "sasthi/saptami", "dvi"), rule!(r"A[nR]Am$", "a", "masc/neut", "sasthi", "bahu"),
        rule!(r"e$", "a", "masc/neut", "saptami", "eka"), rule!(r"ezu$", "a", "masc/neut", "saptami", "bahu"), rule!(r"Ani$", "a", "neut", "prathama/dvitiya", "bahu"),
        rule!(r"iH$", "i", "masc/fem", "prathama", "eka"), rule!(r"im$", "i", "masc/fem/neut", "dvitiya/prathama", "eka"), rule!(r"In$", "i", "masc", "dvitiya", "bahu"),
        rule!(r"IH$", "i", "fem", "prathama/dvitiya", "bahu"), rule!(r"i[nR]A$", "i", "masc/neut", "tritiya", "eka"), rule!(r"yA$", "i", "fem", "tritiya", "eka"),
        rule!(r"aye$", "i", "masc/fem", "caturthi", "eka"), rule!(r"yE$", "i", "fem", "caturthi", "eka"), rule!(r"eH$", "i", "masc/fem", "panchami/sasthi", "eka"),
        rule!(r"yAH$", "i", "fem", "panchami/sasthi", "eka"), rule!(r"yAm$", "i", "fem", "saptami", "eka"), rule!(r"I[nR]Am$", "i", "masc/fem/neut", "sasthi", "bahu"),
        rule!(r"izu$", "i", "masc/fem/neut", "saptami", "bahu"), rule!(r"iByAm$", "i", "masc/fem/neut", "tritiya/caturthi/panchami", "dvi"),
        rule!(r"iByaH$", "i", "masc/fem/neut", "caturthi/panchami", "bahu"), rule!(r"iBiH$", "i", "masc/fem/neut", "tritiya", "bahu"),
        rule!(r"I$", "i", "masc/fem/neut", "prathama/dvitiya", "dvi"), rule!(r"ayaH$", "i", "masc/fem", "prathama", "bahu"), rule!(r"Ini$", "i", "neut", "prathama/dvitiya", "bahu"),
        rule!(r"uH$", "u", "masc/fem", "prathama", "eka"), rule!(r"um$", "u", "masc/fem/neut", "dvitiya/prathama", "eka"), rule!(r"Un$", "u", "masc", "dvitiya", "bahu"),
        rule!(r"UH$", "u", "fem", "prathama/dvitiya", "bahu"), rule!(r"u[nR]A$", "u", "masc/neut", "tritiya", "eka"), rule!(r"vA$", "u", "fem", "tritiya", "eka"),
        rule!(r"ave$", "u", "masc/fem", "caturthi", "eka"), rule!(r"vE$", "u", "fem", "caturthi", "eka"), rule!(r"oH$", "u", "masc/fem", "panchami/sasthi", "eka"),
        rule!(r"vAH$", "u", "fem", "panchami/sasthi", "eka"), rule!(r"vAm$", "u", "fem", "saptami", "eka"), rule!(r"U[nR]Am$", "u", "masc/fem/neut", "sasthi", "bahu"),
        rule!(r"uzu$", "u", "masc/fem/neut", "saptami", "bahu"), rule!(r"uByAm$", "u", "masc/fem/neut", "tritiya/caturthi/panchami", "dvi"),
        rule!(r"uByaH$", "u", "masc/fem/neut", "caturthi/panchami", "bahu"), rule!(r"uBiH$", "u", "masc/fem/neut", "tritiya", "bahu"),
        rule!(r"U$", "u", "masc/fem/neut", "prathama/dvitiya", "dvi"), rule!(r"avaH$", "u", "masc/fem", "prathama", "bahu"), rule!(r"Uni$", "u", "neut", "prathama/dvitiya", "bahu"),
        rule!(r"tA$", "tf", "masc/fem", "prathama", "eka"), rule!(r"tArO$", "tf", "masc/fem", "prathama/dvitiya", "dvi"), rule!(r"tAraH$", "tf", "masc/fem", "prathama", "bahu"),
        rule!(r"tAram$", "tf", "masc/fem", "dvitiya", "eka"), rule!(r"tFn$", "tf", "masc", "dvitiya", "bahu"), rule!(r"tFH$", "tf", "fem", "dvitiya", "bahu"),
        rule!(r"trA$", "tf", "masc/fem", "tritiya", "eka"), rule!(r"tre$", "tf", "masc/fem", "caturthi", "eka"), rule!(r"tuH$", "tf", "masc/fem", "panchami/sasthi", "eka"),
        rule!(r"tari$", "tf", "masc/fem", "saptami", "eka"), rule!(r"tF[nR]Am$", "tf", "masc/fem", "sasthi", "bahu"), rule!(r"tfzu$", "tf", "masc/fem", "saptami", "bahu"),
        rule!(r"tfByAm$", "tf", "masc/fem", "tritiya/caturthi/panchami", "dvi"), rule!(r"tfByaH$", "tf", "masc/fem", "caturthi/panchami", "bahu"), rule!(r"tfBiH$", "tf", "masc/fem", "tritiya", "bahu"),
        
        // HALANTA: 'at' STEMS (Missing previously!)
        rule!(r"an$", "at", "masc", "prathama/sambodhana", "eka"), rule!(r"antO$", "at", "masc", "prathama/dvitiya", "dvi"), rule!(r"antaH$", "at", "masc", "prathama", "bahu"),
        rule!(r"antam$", "at", "masc", "dvitiya", "eka"), rule!(r"ataH$", "at", "masc/neut", "dvitiya/panchami/sasthi", "bahu/eka"),
        rule!(r"atA$", "at", "masc/neut", "tritiya", "eka"), rule!(r"ate$", "at", "masc/neut", "caturthi", "eka"),
        rule!(r"ati$", "at", "masc/neut", "saptami", "eka"), rule!(r"atoH$", "at", "masc/neut", "sasthi/saptami", "dvi"),
        rule!(r"atAm$", "at", "masc/neut", "sasthi", "bahu"), rule!(r"atsu$", "at", "masc/neut", "saptami", "bahu"),
        rule!(r"adByAm$", "at", "masc/neut", "tritiya/caturthi/panchami", "dvi"), rule!(r"adBiH$", "at", "masc/neut", "tritiya", "bahu"), rule!(r"adByaH$", "at", "masc/neut", "caturthi/panchami", "bahu"),
        rule!(r"at$", "at", "neut", "prathama/dvitiya", "eka"), rule!(r"atI$", "at", "neut", "prathama/dvitiya", "dvi"), rule!(r"anti$", "at", "neut", "prathama/dvitiya", "bahu"),

        // CONSONANTS (Halanta)
        rule!(r"A$", "an", "masc", "prathama", "eka"), rule!(r"AnO$", "an", "masc", "prathama/dvitiya", "dvi"), rule!(r"AnaH$", "an", "masc", "prathama", "bahu"),
        rule!(r"Anam$", "an", "masc", "dvitiya", "eka"), rule!(r"nA$", "an", "masc/neut", "tritiya", "eka"),
        rule!(r"I$", "in", "masc", "prathama", "eka"), rule!(r"inO$", "in", "masc", "prathama/dvitiya", "dvi"), rule!(r"inaH$", "in", "masc", "prathama/dvitiya/panchami/sasthi", "bahu"),
        rule!(r"aH$", "as", "neut", "prathama/dvitiya", "eka"), rule!(r"asI$", "as", "neut", "prathama/dvitiya", "dvi"), rule!(r"AMsi$", "as", "neut", "prathama/dvitiya", "bahu"), rule!(r"asA$", "as", "masc/neut", "tritiya", "eka"),
        rule!(r"iH$", "is", "neut", "prathama/dvitiya", "eka"), rule!(r"izI$", "is", "neut", "prathama/dvitiya", "dvi"), rule!(r"IMzi$", "is", "neut", "prathama/dvitiya", "bahu"), rule!(r"izA$", "is", "masc/neut", "tritiya", "eka"), rule!(r"irBiH$", "is", "masc/neut", "tritiya", "bahu"),
        rule!(r"uH$", "us", "neut", "prathama/dvitiya", "eka"), rule!(r"uzI$", "us", "neut", "prathama/dvitiya", "dvi"), rule!(r"UMzi$", "us", "neut", "prathama/dvitiya", "bahu"), rule!(r"uzA$", "us", "masc/neut", "tritiya", "eka"), rule!(r"urBiH$", "us", "masc/neut", "tritiya", "bahu"),
        rule!(r"vAn$", "vat", "masc", "prathama", "eka"), rule!(r"vantO$", "vat", "masc", "prathama/dvitiya", "dvi"), rule!(r"vantaH$", "vat", "masc", "prathama", "bahu"), rule!(r"vatA$", "vat", "masc/neut", "tritiya", "eka"), rule!(r"van$", "vat", "masc", "sambodhana", "eka"),
        rule!(r"mAn$", "mat", "masc", "prathama", "eka"), rule!(r"mantO$", "mat", "masc", "prathama/dvitiya", "dvi"), rule!(r"mantaH$", "mat", "masc", "prathama", "bahu"), rule!(r"matA$", "mat", "masc/neut", "tritiya", "eka"), rule!(r"man$", "mat", "masc", "sambodhana", "eka"),
    ];
}

pub fn get_stems(word: &str) -> Vec<GuessedStem> {
    let mut guessed = Vec::new();
    for rule in STEMMING_RULES.iter() {
        if rule.pattern.is_match(word).unwrap_or(false) {
            let replaced = rule.pattern.replace(word, rule.replacement).into_owned();
            guessed.push(GuessedStem {
                stem: replaced,
                gender: rule.gender.to_string(),
                case: rule.case.to_string(),
                vacana: rule.vacana.to_string(),
            });
        }
    }
    guessed
}
