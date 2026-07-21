import { useEffect, useRef } from "react";

import { EditorView, basicSetup } from "codemirror";
import { EditorState } from "@codemirror/state";

export const DEFAULT_CODE = `// do not trust the terminal
`;

type EditorProps = {
    resetVersion: number;
};

const uplinkTheme = EditorView.theme(
    {
        "&": {
            color: "#d8e4f4",
            backgroundColor: "transparent",
            height: "100%",
        },
        ".cm-content": {
            caretColor: "#6ca5ff",
        },
        ".cm-scroller": {
            fontFamily: '"IBM Plex Mono", "SFMono-Regular", Consolas, monospace',
        },
        ".cm-gutters": {
            backgroundColor: "rgba(255, 255, 255, 0.01)",
            color: "#8ea0bb",
            border: "none",
        },
        ".cm-activeLine": {
            backgroundColor: "rgba(108, 165, 255, 0.06)",
        },
        ".cm-cursor, .cm-dropCursor": {
            borderLeftColor: "#6ca5ff",
        },
        "&.cm-focused .cm-selectionBackground, .cm-selectionBackground, ::selection": {
            backgroundColor: "rgba(108, 165, 255, 0.18) !important",
        },
        ".cm-panels": {
            backgroundColor: "#08101f",
            color: "#d8e4f4",
        },
        ".cm-tooltip": {
            border: "1px solid rgba(108, 165, 255, 0.18)",
            backgroundColor: "#08101f",
            color: "#d8e4f4",
        },
    },
    { dark: true },
);

export default function Editor({ resetVersion }: EditorProps) {
    const editorRef = useRef<HTMLDivElement>(null);
    const viewRef = useRef<EditorView | null>(null);

    useEffect(() => {
        if (!editorRef.current) {
            return;
        }

        const state = EditorState.create({
            doc: DEFAULT_CODE,
            extensions: [basicSetup, uplinkTheme, EditorView.lineWrapping],
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
                insert: DEFAULT_CODE,
            },
        });
    }, [resetVersion]);

    return (
        <section className="panel editor-panel">
            <div className="panel-head">
                <h2 className="panel-title">Console</h2>
                <span className="panel-badge">phosphor link</span>
            </div>
            <div className="panel-body">
                <div ref={editorRef} />
            </div>
        </section>
    );
}