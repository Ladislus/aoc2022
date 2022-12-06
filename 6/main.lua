-- #arg dosen't return the size of the array, but the maximum index, so size - 1
if #arg ~= 1 then io.stderr:write(string.format("Usage: %s input_file_path\n", arg[0])) end

file = io.open(arg[1], 'r')
lines = file:lines()

for line in lines do

    from = 1
    to = 4

    repeat

        cur = string.sub(line, from, to)

        for char in cur:gmatch"." do
            _, count = string.gsub(cur, char, "")
            if count > 1 then goto next_sbustring end
        end

        print(string.format("Found substring \"%s\" without duplicate at [%d, %d]", cur, from, to))
        goto next_input

        ::next_sbustring::

        from = from + 1
        to = to + 1

    until to == (#line + 1)

    ::next_input::

end

io.close(file)