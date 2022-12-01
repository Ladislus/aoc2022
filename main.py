from sys import argv, exit

def help():
    print(f'python3 {argv[0]} input_file_path')
    exit(1)

def read_file(filepath: str) -> list[int]:

    elves: list[int] = []
    calories: int = 0

    with open(filepath) as file:
        for line in file.readlines():
            if line == '\n':
                if calories != 0:
                    elves.append(calories)
                    calories = 0
            else:
                calories += int(line)

    return elves


if __name__ == '__main__':

    if len(argv) != 2:
        help()

    elves: list[int] = read_file(argv[1])
    print(sum(sorted(elves, reverse=True)[0:3]))
