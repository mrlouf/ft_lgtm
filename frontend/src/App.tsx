import Editor from "./components/Editor";
import Output from "./components/Output";
import RunButton from "./components/RunButton";
import StatusBar from "./components/StatusBar";

export default function App() {

    return (

        <main className="container">

            <h1>LGTM Playground</h1>

            <Editor />

            <RunButton />

            <Output />

            <StatusBar />

        </main>

    );

}
