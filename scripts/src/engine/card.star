def _new(
    id="",
    speaker="",
    text="",
    accept_text="",
    reject_text="",
    image_url="",
    on_show=[],
    on_reject=[],
    on_accept=[],
):
    return struct(
        id=id,
        speaker=speaker,
        text=text,
        accept_text=accept_text,
        reject_text=reject_text,
        image_url=image_url,
        on_show=on_show,
        on_reject=on_reject,
        on_accept=on_accept,
    )


card = struct(new=_new)
