import { useEffect, useRef } from "react";

import { EditorView, basicSetup } from "codemirror";
import { EditorState } from "@codemirror/state";

type EditorProps = {
    code: string;
    onChange: (value: string) => void;
    resetVersion: number;
};

export default function Editor({ code, onChange, resetVersion }: EditorProps) {
    const editorRef = useRef<HTMLDivElement>(null);
    const viewRef = useRef<EditorView | null>(null);

    useEffect(() => {
        if (!editorRef.current) {
            return;
        }

        const state = EditorState.create({
            doc: code,
            extensions: [
                basicSetup,
                EditorView.lineWrapping,
                EditorView.updateListener.of((update) => {
                    if (update.docChanged) {
                        onChange(update.state.doc.toString());
                    }
                }),
            ],
        });

        const view = new EditorView({
            state,
            parent: editorRef.current,
        });

        viewRef.current = view;

        return () => {
            view.destroy();
            viewRef.current = null;
        };
    }, []);

    useEffect(() => {
        const view = viewRef.current;
        if (!view) {
            return;
        }

        view.dispatch({
            changes: {
                from: 0,
                to: view.state.doc.length,
                insert: code,
            },
        });
    }, [resetVersion]);

    return (
        <section className="panel editor-panel">
            <div className="panel-head">
                <h2 className="panel-title">Console</h2>
                <span className="panel-badge">Javascript</span>
            </div>
            <div className="panel-body">
                <div ref={editorRef} />
            </div>
        </section>
    );
}