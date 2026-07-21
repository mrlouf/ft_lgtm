export default function RunButton() {
    function handleRun() {
        console.log("Running code...");
    }

    return (
        <div className="controls">
            <button className="run-button" onClick={handleRun}>
                Run Code
            </button>
        </div>
    );
}