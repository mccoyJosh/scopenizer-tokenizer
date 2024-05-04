package tokenizerFactory

import (
	"tp/src/tokenizer"
	"tp/src/tokenizer/tokenizerFactory/javaTokenizer"
	"tp/src/util"
)

func TokenizeGivenLanguageAndFile(language string, path string) []tokenizer.Token {
	tknzr := GetTokenizer(language)

	if tknzr != nil {
		text := util.GetContentOfFile(path)
		return tknzr.Tokenize(text)
	}

	util.Error("Given language for tokenizer is not supported: "+language, nil)
	return nil
}

func GetTokenizer(language string) *tokenizer.Tokenizer {
	var tknzr *tokenizer.Tokenizer = nil

	switch language {
	case "java":
		{
			tknzr = javaTokenizer.Tokenizer()
		}
	case "cpp":
		{
			util.Warning(language + " is not ready to be tokenized yet!")
		}
	case "python":
		{
			util.Warning(language + " is not ready to be tokenized yet!")
		}
	case "fsharp":
		{
			util.Warning(language + " is not ready to be tokenized yet!")
		}
	case "csharp":
		{
			util.Warning(language + " is not ready to be tokenized yet!")
		}
	}

	if tknzr == nil {
		util.Fatal("Given language for tokenizer is not supported: " + language)
	}
	return tknzr
}
