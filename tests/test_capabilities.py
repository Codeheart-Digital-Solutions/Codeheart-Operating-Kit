from codeheart_operating_kit.capabilities import parse_capability_status


def test_native_capability_status_parsing():
    parsed = parse_capability_status("Documents Browser PDF")
    assert parsed["documents"] == "available"
    assert parsed["browser"] == "available"
    assert parsed["pdf"] == "available"
    assert parsed["spreadsheets"] == "unknown"
