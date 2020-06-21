package canvas

// This code is inspired by https://gist.github.com/bumcru/729632c7587f16c69d40a878c0bde750.

var (
	spaceCharTable = map[rune]interface{}{
		'\n': nil, '\t': nil, '\v': nil, ' ': nil, '　': nil,
	}
	startBracketTable = map[rune]interface{}{
		'(': nil, '{': nil, '[': nil, '<': nil,
		'「': nil, '『': nil, '（': nil, '｛': nil, '【': nil, '＜': nil, '≪': nil, '［': nil,
	}
	endCharTable = map[rune]interface{}{
		'、': nil, '。': nil, '.': nil, ',': nil,

		'ぁ': nil, 'ぃ': nil, 'ぅ': nil, 'ぇ': nil, 'ぉ': nil,
		'っ': nil, 'ゃ': nil, 'ゅ': nil, 'ょ': nil,
		'ァ': nil, 'ィ': nil, 'ゥ': nil, 'ェ': nil, 'ォ': nil,
		'ッ': nil, 'ャ': nil, 'ュ': nil, 'ョ': nil,
		'ｧ': nil, 'ｨ': nil, 'ｩ': nil, 'ｪ': nil, 'ｫ': nil,
		'ｯ': nil, 'ｬ': nil, 'ｭ': nil, 'ｮ': nil,

		')': nil, '}': nil, ']': nil, '>': nil,
		'」': nil, '』': nil, '）': nil, '｝': nil, '】': nil, '＞': nil, '≫': nil, '］': nil,

		'・': nil, 'ー': nil, '―': nil, '-': nil,
		'：': nil, '；': nil, '／': nil, '/': nil,
		'ゝ': nil, '々': nil, '！': nil, '？': nil, '!': nil, '?': nil,
	}
)

func spaceChar(r rune) bool {
	_, ok := spaceCharTable[r]
	return ok
}

func oneByteChar(r rune) bool {
	return len(string(r)) == 1
}

func startBracket(r rune) bool {
	_, ok := startBracketTable[r]
	return ok
}

func endChar(r rune) bool {
	_, ok := endCharTable[r]
	return ok
}
