import { useEffect, useRef } from "react";

import { EditorView, basicSetup } from "codemirror";
import { EditorState } from "@codemirror/state";

import LanguageButton from "./LanguageButton";

type Language = "javascript" | "python" | "go";

type EditorProps = {
    code: string;
    language: Language;
    onChange: (value: string) => void;
    onChangeLanguage: (language: Language) => void;
    resetVersion: number;
};

export default function Editor({
    code,
    language,
    onChange,
    onChangeLanguage,
    resetVersion,
}: EditorProps) {
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
    }, [resetVersion, code]);

    return (
        <section className="panel editor-panel">
            <div className="panel-head">
                <h2 className="panel-title">Console</h2>
                <LanguageButton
                    language={language}
                    onChangeLanguage={onChangeLanguage}
                />
            </div>
            <div className="panel-body">
                <div ref={editorRef} />
            </div>
        </section>
    );
}