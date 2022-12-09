import { argv, argv0, exit, listenerCount } from "process";
import { createReadStream } from "fs";
import { createInterface } from "readline";
import { assert } from "console";

class Tree {
    height: number;

    visibleLeft: boolean | null;
    maxLeft: number;
    viewLeft: number;
    visibleTop: boolean | null;
    maxTop: number;
    viewTop: number;
    visibleRight: boolean | null;
    maxRight: number;
    viewRight: number;
    visibleBottom: boolean | null;
    maxBottom: number;
    viewBottom: number;

    constructor(height: number) {
        this.height = height;

        this.visibleLeft = null;
        this.maxLeft = -1;
        this.visibleTop = null;
        this.maxTop = -1;
        this.visibleRight = null;
        this.maxRight = -1;
        this.visibleBottom = null;
        this.maxBottom = -1;
    }

    scenicScore(leftTrees: Tree[], topTrees: Tree[], rightTrees: Tree[], bottomTrees: Tree[]): number {

        // Count trees visible on the left
        let leftScenicScore = 0;
        for (let index = 0; index < leftTrees.length; ++index) {
            ++leftScenicScore;
            if (leftTrees[index].height >= this.height) break;
        }

        // Count trees visible on the top (array was filled in reverse)
        let topScenicScore = 0;
        for (let index = 0; index < topTrees.length; ++index) {
            ++topScenicScore;
            if (topTrees[index].height >= this.height) break;
        }

        // Count trees visible on the right
        let rightScenicScore = 0;
        for (let index = 0; index < rightTrees.length; ++index) {
            ++rightScenicScore;
            if (rightTrees[index].height >= this.height) break;
        }

        // Count trees visible on the bottom
        let bottomScenicScore = 0;
        for (let index = 0; index < bottomTrees.length; ++index) {
            ++bottomScenicScore;
            if (bottomTrees[index].height >= this.height) break;
        }

        return leftScenicScore * topScenicScore * rightScenicScore * bottomScenicScore;
    }
}

class Forest {
    length: number;
    height: number;

    trees: Tree[][];

    constructor(length: number) {
        this.height = 0;
        this.length = length;

        this.trees = [];
    }

    addRow() {
        this.trees.push([]);
        ++this.height;
    }

    addTree(tree: Tree) {

        const lastLine = this.trees[this.height - 1];
        lastLine.push(tree);
    }

    compute() {
        // Compute Tree Left & Top visibility
        this.trees.forEach((row: Tree[], rowIndex: number) => {
            row.forEach((tree: Tree, columnIndex: number) => {
                if (columnIndex == 0) {
                    tree.visibleLeft = true;
                    tree.maxLeft = tree.height;
                } else {
                    const leftTree: Tree = row[columnIndex - 1];
                    tree.visibleLeft = (leftTree.maxLeft < tree.height);
                    tree.maxLeft = Math.max(leftTree.maxLeft, tree.height);
                }

                if (rowIndex == 0) {
                    tree.visibleTop = true;
                    tree.maxTop = tree.height;
                } else {
                    const topTree: Tree = this.trees[rowIndex - 1][columnIndex];
                    tree.visibleTop = (topTree.maxTop < tree.height);
                    tree.maxTop = Math.max(topTree.maxTop, tree.height);
                }
            })
        });

        // Compute Tree Right & Bottom visibility
        for (let rowIndex = this.height - 1; rowIndex >= 0; --rowIndex) {
            for (let columnIndex = this.length - 1; columnIndex >= 0; --columnIndex) {

                const tree: Tree = this.trees[rowIndex][columnIndex];

                if (columnIndex == this.length - 1) {
                    tree.visibleRight = true;
                    tree.maxRight = tree.height;
                } else {
                    const rightTree: Tree = this.trees[rowIndex][columnIndex + 1];
                    tree.visibleRight = (rightTree.maxRight < tree.height);
                    tree.maxRight = Math.max(rightTree.maxRight, tree.height);
                }

                if (rowIndex == this.height - 1) {
                    tree.visibleBottom = true;
                    tree.maxBottom = tree.height;
                } else {
                    const bottomTree: Tree = this.trees[rowIndex + 1][columnIndex];
                    tree.visibleBottom = (bottomTree.maxBottom < tree.height);
                    tree.maxBottom = Math.max(bottomTree.maxBottom, tree.height);
                }
            }
        }

        let count: number = 0;

        // Check if all values where computed
        this.trees.forEach((row: Tree[], rowIndex: number) => {
            row.forEach((tree: Tree, columnIndex: number) => {
                assert(tree.visibleLeft != null, `Tree [${rowIndex}:${columnIndex}] has null left visibility`);
                assert(tree.visibleTop != null, `Tree [${rowIndex}:${columnIndex}] has null top visibility`);
                assert(tree.visibleRight != null, `Tree [${rowIndex}:${columnIndex}] has null right visibility`);
                assert(tree.visibleBottom != null, `Tree [${rowIndex}:${columnIndex}] has null bottom visibility`);

                if (tree.visibleLeft || tree.visibleTop || tree.visibleRight || tree.visibleBottom) ++count
            })
        });

        console.log(`input v1: ${count}`);
    }
}

function help() {
    console.error(`Usage: ${argv0} input_file_path`)
    exit(1)
}

async function input_v1(filepath: string) {

    let forest: Forest = null;

    // How hard can it be to open a simple file, seriously ?
    const lineReader = createInterface({
        input: createReadStream(filepath),
        crlfDelay: Infinity
    });

    for await (const line of lineReader) {

        const splittedLine = line.split('');

        if (forest == null) forest = new Forest(splittedLine.length);
        forest.addRow();

        splittedLine.forEach((value: string, _: number) => forest.addTree(new Tree(parseInt(value))));
    }

    forest.compute();
}

async function input_v2(filepath: string) {

    let forest: Forest = null;

    // How hard can it be to open a simple file, seriously ?
    const lineReader = createInterface({
        input: createReadStream(filepath),
        crlfDelay: Infinity
    });

    for await (const line of lineReader) {

        const splittedLine = line.split('');

        if (forest == null) forest = new Forest(splittedLine.length);
        forest.addRow();

        splittedLine.forEach((value: string, _: number) => forest.addTree(new Tree(parseInt(value))));
    }

    let maxScenicScore: number = -1;

    forest.trees.forEach((row: Tree[], rowIndex: number) => {
        row.forEach((tree: Tree, columnIndex: number) => {
            let leftTrees: Tree[] = [];
            let topTrees: Tree[] = [];
            let rightTrees: Tree[] = [];
            let bottomTrees: Tree[] = [];

            if (columnIndex > 0) leftTrees = forest.trees[rowIndex].slice(0, columnIndex).reverse();
            if (columnIndex < forest.length - 1) rightTrees = forest.trees[rowIndex].slice(columnIndex + 1);
            if (rowIndex > 0) {
                for (let currentRownIndex = rowIndex - 1; currentRownIndex >= 0; --currentRownIndex) topTrees.push(forest.trees[currentRownIndex][columnIndex]);
            }
            if (rowIndex < forest.height - 1) {
                for (let currentRowIndex = rowIndex + 1; currentRowIndex <= forest.length - 1; ++currentRowIndex) bottomTrees.push(forest.trees[currentRowIndex][columnIndex]);
            }

            let currentScenicScore: number = tree.scenicScore(leftTrees, topTrees, rightTrees, bottomTrees);

            maxScenicScore = Math.max(maxScenicScore, currentScenicScore);
        })
    });

    console.log(`input_v2: ${maxScenicScore}`);
}



if (argv.length != 3) help();
input_v1(argv[2]);
input_v2(argv[2]);