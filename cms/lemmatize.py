#!/usr/bin/env python3
"""Лемматизация текстового файла с помощью pymorphy3.

Использование:
    python lemmatize_file.py текст.txt
    python lemmatize_file.py текст.txt -o леммы.txt
    python lemmatize_file.py текст.txt --counts --top 50
"""

import argparse
import re
import sys
from collections import Counter
from functools import lru_cache

import pymorphy3

morph = pymorphy3.MorphAnalyzer()

TOKEN_RE = re.compile(r'\w+')


@lru_cache(maxsize=200_000)
def normal_form(token):
    """Нормальная форма одного токена. Кэш даёт ускорение в разы:
    в реальном тексте одни и те же слова встречаются постоянно."""
    return morph.parse(token)[0].normal_form


def lemmatize(text, fix_yo=True):
    """Возвращает список лемм для строки текста."""
    if fix_yo:
        text = text.replace('ё', 'е').replace('Ё', 'Е')
    tokens = TOKEN_RE.findall(text.lower())
    return [normal_form(t) for t in tokens]


def iter_lemmas(path, encoding='utf-8', fix_yo=True):
    """Читает файл построчно и отдаёт леммы по одной.

    Построчное чтение вместо f.read() — чтобы файл на несколько гигабайт
    не пришлось целиком держать в памяти."""
    with open(path, encoding=encoding, errors='replace') as f:
        for line in f:
            yield from lemmatize(line, fix_yo=fix_yo)


def main():
    parser = argparse.ArgumentParser(description='Лемматизация текстового файла')
    parser.add_argument('filename', help='путь к входному файлу')
    parser.add_argument('-o', '--output', help='файл для записи результата (по умолчанию stdout)')
    parser.add_argument('-e', '--encoding', default='utf-8',
                        help='кодировка входного файла (utf-8, cp1251, ...)')
    parser.add_argument('--counts', action='store_true',
                        help='вывести частотный список вместо потока лемм')
    parser.add_argument('--top', type=int, default=None,
                        help='сколько самых частых лемм показать (с --counts)')
    parser.add_argument('--keep-yo', action='store_true',
                        help='не заменять ё на е')
    args = parser.parse_args()

    fix_yo = not args.keep_yo

    try:
        lemmas = iter_lemmas(args.filename, args.encoding, fix_yo)
        out = open(args.output, 'w', encoding='utf-8') if args.output else sys.stdout

        if args.counts:
            counter = Counter(lemmas)
            for lemma, n in counter.most_common(args.top):
                print(f'{n}\t{lemma}', file=out)
            print(f'Токенов: {sum(counter.values())}, уникальных лемм: {len(counter)}',
                  file=sys.stderr)
        else:
            total = 0
            for lemma in lemmas:
                print(lemma, file=out)
                total += 1
            print(f'Обработано токенов: {total}', file=sys.stderr)

        if args.output:
            out.close()

    except FileNotFoundError:
        sys.exit(f'Файл не найден: {args.filename}')
    except UnicodeDecodeError:
        sys.exit(f'Не удалось прочитать файл в кодировке {args.encoding}. '
                 f'Попробуйте --encoding cp1251')


if __name__ == '__main__':
    main()