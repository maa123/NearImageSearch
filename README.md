# NearImageSearch

[![Actions Status](https://github.com/maa123/NearImageSearch/workflows/Go/badge.svg)](https://github.com/maa123/NearImageSearch/actions)


Goで書いた類似画像検索用のプログラム

動作には事前に別のプログラムで生成したハッシュ値のファイルが必要になります

dHashを利用し、dHashへの変換処理に3ms程度かかり、けもフレ1期のindex(約40万件)の検索に1ms程度かかります
