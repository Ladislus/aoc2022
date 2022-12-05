import std.stdio: writeln, File, stderr;
import core.stdc.stdlib: exit;
import std.exception: ErrnoException;
import std.container: DList, Array;
import std.ascii: isWhite;
import std.string: strip, isNumeric;
import std.array: split, array;
import std.range: stride, walkLength;
import std.regex: regex, matchFirst;
import std.conv: to;

void help(string program_name) {
	stderr.writeln("Usage: ", program_name, " input_file_path");
	exit(1);
}

void input_v1(string filepath) {

	auto stacks = Array!(DList!dchar)();

	try {
		File file = File(filepath, "r");

		string current_line;
		while ((current_line = file.readln()) !is null) {

			if (strip(current_line).length == 0) break;

			dchar[] splitted = stride(current_line[1 .. $ - 1].dup, 4).array();

			if (splitted.length == 0) {
				stderr.writeln("Empty stride");
				exit(1);
			}

			// Skip line with numbers
			if (isNumeric(splitted)) continue;

			if (splitted.length > stacks.length)
				foreach (i; stacks.length .. splitted.length)
					stacks.insert(DList!dchar());

			foreach (i; 0..splitted.length)
				if (!isWhite(splitted[i]))
					stacks[i].insertFront(splitted[i]);
		}

		auto reg = regex(r"^move (?P<count>\d+) from (?P<from>\d+) to (?P<to>\d+)$");
		while ((current_line = strip(file.readln())) !is null) {

			auto match = matchFirst(current_line, reg);
			assert(match.length == 4);

			const auto count = (match["count"].to!int());
			const auto from = (match["from"].to!int()) - 1;
			const auto to = (match["to"].to!int()) - 1;

			foreach (i; 0 .. count) {
				assert(from < stacks.length && to < stacks.length && from != to);
				if (walkLength(stacks[from][]) == 0) { writeln("Empty stack, breaking"); break; }

				auto popped = stacks[from].back;
				stacks[from].removeBack();
				stacks[to].insertBack(popped);
			}
		}

		file.close();

		string result = "";
		foreach (DList!dchar key; stacks) {
			if (walkLength(key[]) == 0) result ~= " ";
			else result ~= key.back();
		}

		writeln("result v1: \"", result, "\"");

	} catch	(ErrnoException ex) {
		stderr.writeln("Error while openning file ", filepath, " with error:");
		stderr.writeln(ex.msg);
		exit(1);
	}
}

void input_v2(string filepath) {

	auto stacks = Array!(DList!dchar)();

	try {
		File file = File(filepath, "r");

		string current_line;
		while ((current_line = file.readln()) !is null) {

			if (strip(current_line).length == 0) break;

			dchar[] splitted = stride(current_line[1 .. $ - 1].dup, 4).array();

			if (splitted.length == 0) {
				stderr.writeln("Empty stride");
				exit(1);
			}

			// Skip line with numbers
			if (isNumeric(splitted)) continue;

			if (splitted.length > stacks.length)
				foreach (i; stacks.length .. splitted.length)
					stacks.insert(DList!dchar());

			foreach (i; 0..splitted.length)
				if (!isWhite(splitted[i]))
					stacks[i].insertFront(splitted[i]);
		}

		auto reg = regex(r"^move (?P<count>\d+) from (?P<from>\d+) to (?P<to>\d+)$");
		while ((current_line = strip(file.readln())) !is null) {

			auto match = matchFirst(current_line, reg);
			assert(match.length == 4);

			const auto count = (match["count"].to!int());
			const auto from = (match["from"].to!int()) - 1;
			const auto to = (match["to"].to!int()) - 1;

			dchar[] to_append;

			foreach (i; 0 .. count) {
				assert(from < stacks.length && to < stacks.length && from != to);
				if (walkLength(stacks[from][]) == 0) { writeln("Empty stack, breaking"); break; }

				auto popped = stacks[from].back;
				stacks[from].removeBack();

				to_append ~= popped;
			}

			foreach_reverse (elem; to_append) stacks[to].insertBack(elem);
		}

		file.close();

		string result = "";
		foreach (DList!dchar key; stacks) {
			if (walkLength(key[]) == 0) result ~= " ";
			else result ~= key.back();
		}

		writeln("result v1: \"", result, "\"");

	} catch	(ErrnoException ex) {
		stderr.writeln("Error while openning file ", filepath, " with error:");
		stderr.writeln(ex.msg);
		exit(1);
	}
}

void main(string[] args) {

	if (args.length != 2) help(args[0]);

	input_v1(args[1]);
	input_v2(args[1]);
}
