# WeatherDump

Multi-platform software for record, demodulate, decode and process data from weather satellites.

## Supported Datalink Protocols

| Protocol | Complete Name | Satellites | Band | Support Level |
| -------- | ------------- | ---------- | ---- | ------------- |
| LRPT | Low Rate Picture Transfer | Meteor-MN2 | VHF | Alpha |
| HRD | High Rate Data | NOAA-20 & Suomi | X-Band | Beta |
| APT | Automatic Picture Transfer | NOAA-15, NOAA-18 & NOAA-19 | VHF | Planned (Beta 1) |

## Example Usage

Decoding and processing a Meteor-MN2 soft-symbol file:

```bash
weatherdump lrpt soft ./file_path.bin
```

## Known Bugs

- The LRPT RGB composite is unsynchonized in most occasions. Will be corrected in Beta 1.
- Garbage collection bug causes the app to use a huge amout of memory. Will be corrected in Beta 1.

## Upcoming Features List

The WeatherDump project roadmap is available in our [Notion Page](https://www.notion.so/fef088dd80b34bd9a6547e890ed962d3?v=182d497ce07d42bc98fb325ca091e813).

## Changelog

### Alpha 2

- Support for multiple decoders for each protocol.
- Added support to scrambled unsynchronized CADU files for HRD.
- Added support to synchronized unscrambled CADU files for HRD.
- Fix multi-thread image processing freezing.
- New decoder and processor progress indicator.
- New CLI argument for multiple decoders.
- Improved far from perfect documentation.
- First public release of GUI version.
- Statistical SCID recover.
- New LRPT frame stacker with proper line synchronization.
- Added RGB multispectral composites for LRPT.
- Exported functions better documentated.
- GUI stylesheet refactor following Airbnb's styleguide.
- Webfonts now loaded from the CSS.
- Improving multi-theme support.
- Implemented Golang module manager.
- Tabs routing to the right place.
- New Javascript library for the REST API.
- Fix WebSockets synchronization.
- General improvements to the REST API.
- Unified WebSockets handlers.
- New thumbnail generator API.
- Add support for multiple apps opened at the same time.
- New engine handler is now running on client-side.