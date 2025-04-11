// # Dir Kit
//
// * A simple toolkit for folder and file handling that eliminates
// boilerplate or wraps commonly used functions in a consistent
// namespace for easy rememberance/importing.

using System;
using System.Collections.Generic;
using System.IO;


namespace dirkit
{
    public static class FileUtils
    {
        public static List<string>? GetDirContents(string path, bool fullPath = false)
        {
            if (System.IO.Directory.Exists(path))
            {
                string[] entries = System.IO.Directory.GetFileSystemEntries(path);
                List<string> result = new List<string>();

                foreach (string entry in entries)
                {
                    if (fullPath)
                    {
                        result.Add(entry);
                    }
                    else
                    {
                        result.Add(System.IO.Path.GetFileName(entry));
                    }
                }

                return result;
            }

            return null;
        }
    }
}

