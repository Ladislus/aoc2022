const std = @import("std");

fn help(program_name: [:0]const u8) noreturn {
    std.log.err("Usage is \"{s} input_file_path\"\n", .{ program_name });
    std.os.exit(1);
}

fn openAndReadFile(filepath: [:0]const u8) !void {
    var score: u32 = 0;

    const file = try std.fs.cwd().openFile(filepath, .{ .mode = std.fs.File.OpenMode.read_only });
    defer file.close();

    var buf: [1024]u8 = undefined;
    while (try file.reader().readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var comp1 = line[0..(line.len / 2)];
        var comp2 = line[(line.len / 2)..];

        std.sort.sort(u8, comp1, {}, comptime std.sort.asc(u8));
        std.sort.sort(u8, comp2, {}, comptime std.sort.asc(u8));

        const length: usize = comp1.len;
        var iter1: usize = 0;
        var iter2: usize = 0;

        while ((iter1 < length) and (iter2 < length)) {
            const elem1 = comp1[iter1];
            const elem2 = comp2[iter2];

            if (elem1 == elem2) {
                if (std.ascii.isUpper(elem1)) { score += (elem1 - 65) + 27; }
                else { score += elem1 - 96; }
                break;
             }
            else if (elem1 < elem2) { iter1 += 1; }
            else { iter2 += 1; }
        }
    }

    std.log.info("Score: {}", .{ score });
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{ .safety = true }) {};
    defer {
        const leaked = gpa.deinit();
        if (leaked) @panic("GeneralPurposeAllocator leaked !");
    }
    const allocator = gpa.allocator();
    const argv = try std.process.argsAlloc(allocator);
    defer std.process.argsFree(allocator, argv);

    if (argv.len != 2) help(argv[0]);

    try openAndReadFile(argv[1]);
}
