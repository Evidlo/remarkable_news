package main

import(
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func pdf_img(url string, output string) {
    os.Setenv("PATH", "/opt/bin")
    magickExecutable, _ := exec.LookPath("convert")
    debug("Coverting pdf ", url)

    builder := strings.Builder{}
    builder.WriteString("png:")
    builder.WriteString(output)
    outpng := builder.String();

    cmd := &exec.Cmd {
        Path: magickExecutable,
        Args: []string{ magickExecutable, 
            "-density", "226", 
            url, 
            "-resize","1404",
            "-colorspace", "gray",
            "-gravity", "north",
            "-crop", "1404x1872+0+100",
            "+repage",
            outpng,
        },
        Stdout: os.Stdout,
        Stderr: os.Stdout,
    }

    if err := cmd.Run(); err != nil {
        fmt.Println("Error:", err)
    }
}