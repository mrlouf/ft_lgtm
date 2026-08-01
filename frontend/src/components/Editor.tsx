import { useEffect, useRef } from "react";

import { EditorView, basicSetup } from "codemirror";
import { EditorState } from "@codemirror/state";

import LanguageButton from "./LanguageButton";

import EnterCIDButton from "./EnterCIDButton";

type Language = "javascript" | "python" | "go";


type EditorProps = {
    code: string;
    language: Language;
    cid: string;
    onChange: (value: string) => void;
    onChangeLanguage: (language: Language) => void;
    onEnterCID: (nextCid: string) => void;
    resetVersion: number;
};

export default function Editor({
    code,
    language,
    cid,
    onChange,
    onChangeLanguage,
    onEnterCID,
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

        const currentValue = view.state.doc.toString();
        if (currentValue === code) {
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

                <EnterCIDButton
                    cid={cid}
                    onEnterCID={onEnterCID}
                />

                <LanguageButton
                    language={language}
                    onChangeLanguage={onChangeLanguage}
                />
            </div>
            <div className="panel-body editor-body">
                <div ref={editorRef} className="editor-mount" />
            </div>
        </section>
    );
}