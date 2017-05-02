words = {
    1: 'gun',
    2: 'shoe',
    3: 'tree',
    4: 'door',
    5: 'hive',
    6: 'trick',
    7: 'devon',
    8: 'plate',
    9: 'line',
    10: 'pen',
    11: 'elevator',
    12: 'shelve',
    13: 'flea',
    14: 'floor',
    15: 'knife',
    16: 'sex',
    17: 'heaven',
    18: 'ate',
    19: 'aisle',
    20: 'plenty',
}

import random


for x in xrange(0, 100):
    word_keys_shuffled = words.keys()
    random.shuffle(word_keys_shuffled)
    for position in word_keys_shuffled:
        guess = raw_input('\nEnter your answer for {}: '.format(position))
        if guess == words[position]:
            print 'yes!\n'
        else:
            print 'no, it was: {}\n'.format(words[position])

