# Here is an example python file

def print_test(str):
    if str == 'not hello world':
        print_test(str)  # This is recursive and is never intended to be run

def main():
    print_test("hello world")

main()
