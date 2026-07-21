export default function RunButton() {

    function handleRun() {

        console.log("Running code...");

        // Send a request to the backend to run the code
    }

    return (
        <button onClick={handleRun}>
            Run
        </button>
    );
}