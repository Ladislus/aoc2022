const std = @import("std");

fn help(program_name: [:0]const u8) noreturn {
    std.log.err("Usage is \"{s} input_file_path\"\n", .{ program_name });
    std.os.exit(1);
}

fn openAndReadFile_v1(filepath: [:0]const u8) !void {
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

    std.log.info("Score v1: {}", .{ score });
}

fn openAndReadFile_v2(filepath: [:0]const u8) !void {
    var score: u32 = 0;

    const file = try std.fs.cwd().openFile(filepath, .{ .mode = std.fs.File.OpenMode.read_only });
    defer file.close();

    var gpa = std.heap.GeneralPurposeAllocator(.{ .safety = true }) {};
    defer {
        const leaked = gpa.deinit();
        if (leaked) @panic("GeneralPurposeAllocator leaked !");
    }
    const allocator = gpa.allocator();
    var common_elements = std.ArrayList(u8).init(allocator);
    defer common_elements.deinit();

    var buf1: [1024]u8 = undefined;
    var buf2: [1024]u8 = undefined;
    var buf3: [1024]u8 = undefined;
    while (true) {

        const tmp1 = try file.reader().readUntilDelimiterOrEof(&buf1, '\n');
        const tmp2 = try file.reader().readUntilDelimiterOrEof(&buf2, '\n');
        const tmp3 = try file.reader().readUntilDelimiterOrEof(&buf3, '\n');

        if ((tmp1 == null) or (tmp2 == null) or (tmp3 == null)) break;

        const elf1 = tmp1.?;
        const elf2 = tmp2.?;
        const elf3 = tmp3.?;

        std.sort.sort(u8, elf1, {}, comptime std.sort.asc(u8));
        std.sort.sort(u8, elf2, {}, comptime std.sort.asc(u8));
        std.sort.sort(u8, elf3, {}, comptime std.sort.asc(u8));

        var iter1: usize = 0;
        var iter2: usize = 0;
        var iter3: usize = 0;

        while ((iter1 < elf1.len) and (iter2 < elf2.len)) {
            const elem1 = elf1[iter1];
            const elem2 = elf2[iter2];

            if (elem1 == elem2) {
                try common_elements.append(elem1);
                iter1 += 1;
                iter2 += 1;
             }
            else if (elem1 < elem2) { iter1 += 1; }
            else { iter2 += 1; }
        }

        iter2 = 0;
        const ce_slice = common_elements.items;

        while ((iter2 < ce_slice.len) and (iter3 < elf3.len)) {
            const elem2 = ce_slice[iter2];
            const elem3 = elf3[iter3];

            if (elem2 == elem3) {
                if (std.ascii.isUpper(elem2)) { score += (elem2 - 65) + 27; }
                else { score += elem2 - 96; }
                common_elements.clearRetainingCapacity();
                break;
             }
            else if (elem2 < elem3) { iter2 += 1; }
            else { iter3 += 1; }
        }
    }

    std.log.info("Score v2: {}", .{ score });
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

    try openAndReadFile_v1(argv[1]);
    try openAndReadFile_v2(argv[1]);
}
