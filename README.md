# NearImageSearch

[![CircleCI](https://circleci.com/gh/maa123/NearImageSearch.svg?style=svg)](https://circleci.com/gh/maa123/NearImageSearch)

[![Actions Status](https://github.com/maa123/NearImageSearch/workflows/Go/badge.svg)](https://github.com/maa123/NearImageSearch/actions)

[![Maintainability](https://api.codeclimate.com/v1/badges/2f023ffbcc864cebb217/maintainability)](https://codeclimate.com/github/maa123/NearImageSearch/maintainability)

Goで書いた類似画像検索用のプログラム

動作には事前に別のプログラムで生成したハッシュ値のファイルが必要になります

dHashを利用し、dHashへの変換処理に3ms程度かかり、けもフレ1期のindex(約40万件)の検索に1ms程度かかります
