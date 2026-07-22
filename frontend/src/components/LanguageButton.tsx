type Language = "javascript" | "python" | "go";

type LanguageButtonProps = {
    language: Language;
    onChangeLanguage: (language: Language) => void;
};

const labels: Record<Language, string> = {
    javascript: "JavaScript",
    python: "Python",
    go: "Go",
};

export default function LanguageButton({
    language,
    onChangeLanguage,
}: LanguageButtonProps) {
    return (
        <select
            className="panel-badge language-select"
            value={language}
            onChange={(event) => onChangeLanguage(event.target.value as Language)}
            aria-label="Select language"
        >
            <option value="javascript">{labels.javascript}</option>
            <option value="python">{labels.python}</option>
            <option value="go">{labels.go}</option>
        </select>
    );
}