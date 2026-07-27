import { useEffect, useState } from "react";

export default function EnterCIDButton({
    cid,
    onEnterCID,
}: {
    cid: string;
    onEnterCID: (nextCid: string) => void;
}) {
    const [inputValue, setInputValue] = useState(cid);

    useEffect(() => {
        setInputValue(cid);
    }, [cid]);

    function handleEnterClick() {
        onEnterCID(inputValue.trim());
    }

    return (
        <div className="enter-cid-button">
            <input
                type="text"
                value={inputValue}
                onChange={(event) => setInputValue(event.target.value)}
                placeholder="Enter CID"
            />
            <button onClick={handleEnterClick}>Enter</button>
        </div>
    );
}