package domainerr

import "errors"

// ErrNotFound はリポジトリ層がエンティティを見つけられなかったときに返す。
// 文字列比較ではなく errors.Is で判定できる。
var ErrNotFound = errors.New("not found")
