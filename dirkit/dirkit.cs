// # Dir Kit
//
// * A simple toolkit for folder and file handling that eliminates
// boilerplate or wraps commonly used functions in a consistent
// namespace for easy rememberance/importing.


using System.Net.Mime;

namespace sandbox
{
    public static class FileUtils
    {
        const string SAFETY_PATH = "D:/safety/";

        public static string GetDate()
        {
            return DateTime.Now.ToString("yyyymmdd");
        }

        public static string GetTime()
        {
            return DateTime.Now.ToString("HH:mm:ss:ff");
        }

        public static List<string>? GetDirContents(string path, bool fullPath = false)
        {
            if (!Directory.Exists(path))
            {
                return null;
            }

            string[] entries = Directory.GetFileSystemEntries(path);
            List<string> result = [];

            foreach (string entry in entries)
            {
                if (fullPath)
                {
                    result.Add(entry);
                }
                else
                {
                    result.Add(Path.GetFileName(entry));
                }
            }

            return result;
        }

        public static void CreateDirectory(string path)
        {
            if (!Directory.Exists(path))
            {
                Directory.CreateDirectory(path);
            }
        }

        public static string CreateDatedDirectory(string path)
        {
            var DatePath = Path.Combine(path, GetDate());
            CreateDirectory(DatePath);
            return DatePath;
        }

        public static void DeleteSafeFile(string path)
        {
            if (!Directory.Exists(path))
            {
                return;
            }

            if (!path.Contains(SAFETY_PATH))
            {
                string safety_path = SAFETY_PATH.ToString();
                string message = $"Path must be within {safety_path}!";
                throw new Exception(message);
            }

            try
            {
                File.Delete(path);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error deleting file {path}: {ex.Message}");
            }
        }

        public static void DeleteSafeFilesInDirectory(string path)
        {
            var Conents = GetDirContents(path, true);
            if (Conents == null){
                return;
            }

            foreach (string file in Conents)
            {
                DeleteSafeFile(file);
            }
        }

        private static void DeleteFolders(string path)
        {
            foreach (string dir in Directory.GetDirectories(path))
            {
                DeleteSafeDirectory(path);

                try
                {
                    Directory.Delete(dir);
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error deleting directory {dir}: {ex.Message}");
                }
            }
        }

        public static void DeleteSafeDirectory(string path)
        {
            if (!Directory.Exists(path))
            {
                return;
            }

            if (!path.Contains(SAFETY_PATH))
            {
                string safety_path = SAFETY_PATH.ToString();
                string message = $"Path must be within {safety_path}!";
                throw new Exception(message);
            }

            var contents = GetDirContents(path, true);
            if (contents == null)
            {
                Directory.Delete(path);
                return;
            }

            DeleteSafeFilesInDirectory(path);
            DeleteFolders(path);

            try
            {
                Directory.Delete(path);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error deleting directory {path}: {ex.Message}");
            }
        }

        public static (string, string, string) GetOsInfo()
        {
            string system = Environment.OSVersion.Platform.ToString();
            string release = Environment.OSVersion.Version.ToString();
            string version = Environment.OSVersion.VersionString;

            return (system, release, version);
        }
    }
}
