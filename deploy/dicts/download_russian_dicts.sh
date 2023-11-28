# Set up russians dicts for postgres full text search

wget https://bitbucket.org/Shaman_Alex/russian-dictionary-hunspell/downloads/ru_RU_UTF-8_20131101.zip
unzip ru_RU_UTF-8_20131101.zip

iconv -o russian.affix ru_RU.aff
iconv -o russian.dict ru_RU.dic

rm -f ru_RU_UTF-8_20131101.zip ru_RU.aff ru_RU.dic
