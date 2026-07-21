type ResetButtonProps = {
    onReset: () => void;
};

export default function ResetButton({ onReset }: ResetButtonProps) {
    return (
        <div className="controls">
            <button className="run-button" onClick={onReset}>
                Reset Code
            </button>
        </div>
    );
}