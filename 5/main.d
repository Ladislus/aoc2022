import std.stdio: writeln, File, stderr;
import core.stdc.stdlib: exit;
import std.exception: ErrnoException;
import std.container: DList, Array;
import std.ascii: newline;
import std.string: strip;
import std.array: split, array, empty, minimallyInitializedArray;
import std.range: stride;

void help(string program_name) {
	stderr.writeln("Usage: ", program_name, " input_file_path");
	exit(1);
}

void input_v1(string filepath) {

	Array!(DList!char) stacks = Array!(DList!char)();
	stacks.insertAfter(DList!char());
	stderr.writefln("Stack size: ", stacks.length);
	exit(1);

	try {
		File file = File(filepath, "r");

		string current_line;
		while ((current_line = strip(file.readln())) !is null) {
			writeln("\"", current_line, "\"");

			if (current_line.length == 0) break;

			dchar[] splitted = stride(current_line[1 .. $ - 1].dup, 4).array();
			
			if (splitted.length == 0) {
				stderr.writeln("Empty stride");
				exit(1);
			}

			if (splitted.length > stacks.length) {
				auto diff = splitted.length - stacks.length;
				// foreach (i; 0..diff) stacks.insertAfter(DList!char());
			}

			writeln("\"", splitted, "\"");
		}

		// while ((current_line = strip(file.readln())) !is null) {
		// 	writeln("\"", current_line, "\"");
		// }

		file.close();

		string result = "";
		foreach (DList!char key; stacks) {
			result ~= key.back();
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
}
