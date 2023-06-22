Console.WriteLine("Hello, what is your name?");

var name = Console.ReadLine();

if (name == string.Empty)
    name = "No Name";

Console.WriteLine("Welcome {0} to your OOP adventure", name);
