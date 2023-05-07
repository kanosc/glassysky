import sys

# sys.argv是system arg variable的拼写.
#sys.argv是一个不断可以添加的参数的列表，这个参数是需要手动在命令行输入的，
print('all arg is :', sys.argv)

if sys.argv[1] == 'zara':
    print('zara is beautiful')
elif sys.argv[2] == 'lihuaiyuan':
    print('my best love is lihuaiyuan!')
else:
    print(f'argv is is else....{sys.argv[3:]}')

