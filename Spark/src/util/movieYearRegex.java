package util;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class movieYearRegex {
	private static String moduleType = ".* \\(([1-9][0-9][0-9][0-9])\\).*";

	public static void main(String[] args) {
		System.out.println(movieYearReg("GoldenEye (1995)"));
	}

	public static int movieYearReg(String str) {
		int retYear = 1994;
		Pattern patternType = Pattern.compile(moduleType);
		Matcher matcherType = patternType.matcher(str);
		while (matcherType.find()) {
			retYear = Integer.parseInt(matcherType.group(1));
		}
		return retYear;
	}
}
