import { argv, argv0, exit } from "process";
import { createReadStream } from "fs";
import { createInterface } from "readline";

class Tree {
    visible: { [key: string]: boolean | null };
    height: number;

    constructor(height: number) {
        this.height = height;
        this.visible = { "left": null, "top": null, "right": null, "bottom": null };
    }
}

function help() {
    console.error(`Usage: ${argv0} input_file_path`)
    exit(1)
}

async function input(filepath: string) {

    let elements: string[][] = [];

    // How hard can it be to open a simple file, seriously ?
    const lineReader = createInterface({
        input: createReadStream(filepath),
        crlfDelay: Infinity
    });

    for await (const line of lineReader) elements.push(line.split(""));


}

if (argv.length != 3) help();
input(argv[2]);