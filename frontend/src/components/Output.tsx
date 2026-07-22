type OutputProps = {
    output: string;
};

export default function Output({ output }: OutputProps) {
    return (
        <section className="panel output-panel">
            <div className="panel-head">
                <h2 className="panel-title">Output</h2>
                <span className="panel-badge">live feed</span>
            </div>
            <div className="panel-body">
                <pre>{output}</pre>
            </div>
        </section>
    );
}