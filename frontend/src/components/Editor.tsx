import { useEffect, useRef } from "react";

import { EditorView, basicSetup } from "codemirror";
import { EditorState } from "@codemirror/state";


export default function Editor() {

    const editorRef = useRef<HTMLDivElement>(null);

    useEffect(() => {

        if (!editorRef.current) {
            return;
        }

        const state = EditorState.create({
            doc: "// LGTM Playground\n",
            extensions: [
                basicSetup,
            ],
        });


        const view = new EditorView({
            state,
            parent: editorRef.current,
        });


        return () => {
            view.destroy();
        };

    }, []);


    return (
        <section className="editor">

            <h2>Editor</h2>

            <div ref={editorRef} />

        </section>
    );
}