import subprocess
import os

def run(args):
    if not args:
        print("Usage: cp-cli test filename.extension")
        return

    filename = args[0]
    extension = filename.split(".")[-1]

    numTests = int(input("Enter number of testcases: "))
    testcases = []

    for i in range(numTests):
        print(f"\nTest case #{i+1}")
        testInput = input("Input:\n")
        expectedOutput = input("Expected Output:\n")
        testcases.append((testInput, expectedOutput))

    print("\nRunning tests...\n")

    exeName = "a.out"

    if extension in ["cpp", "cc"]:
        subprocess.run(["g++", filename, "-o", exeName])
        runCmd = f"./{exeName}"
    elif extension == "c":
        subprocess.run(["gcc", filename, "-o", exeName])
        runCmd = f"./{exeName}"
    elif extension == "py":
        runCmd = ["python3", filename]
    elif extension == "java":
        subprocess.run(["javac", filename])
        runCmd = ["java", filename.replace(".java", "")]
    else:
        print(f"Unsupported extension: .{extension}")
        return

    for i, (inputData, expectedOutput) in enumerate(testcases):
        print(f"\n[TEST {i+1}]")
        try:
            if isinstance(runCmd, str):
                result = subprocess.run(runCmd.split(), input=inputData, text=True, capture_output=True, timeout=5)
            else:
                result = subprocess.run(runCmd, input=inputData, text=True, capture_output=True, timeout=5)

            output = result.stdout.strip()
            expected = expectedOutput.strip()

            if output == expected:
                print("✅ Passed")
            else:
                print("❌ Failed")
                print(f"Expected: {expected}")
                print(f"Got     : {output}")
        except Exception as e:
            print(f"❌ Error running test: {e}")

    if os.path.exists(exeName):
        os.remove(exeName)