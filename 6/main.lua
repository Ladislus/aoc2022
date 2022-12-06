function input_v1(filepath)

    file = io.open(filepath, 'r')
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

end

function input_v2(filepath)

    file = io.open(filepath, 'r')
    lines = file:lines()

    for line in lines do

        from = 1
        to = 14

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

end

if #arg ~= 1 then io.stderr:write(string.format("Usage: %s input_file_path\n", arg[0])) end
input_v1(arg[1])
input_v2(arg[1])