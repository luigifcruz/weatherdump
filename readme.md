# WeatherDump

### Supported Datalink Protocols
| Protocol | Complete Name | Satellites | Band | Support Level |
| -------- | ------------- | ---------- | ---- | ------------- |
| LRPT | Low Rate Picture Transfer | Meteor-MN2 | VHF | Alpha |
| HRD | High Rate Data | NOAA-20 & Suomi | X-Band | Beta |
| APT | Automatic Picture Transfer | NOAA-15, NOAA-18 & NOAA-19 | VHF | Planned (Beta 2) |

### Example Usage
Decoding and processing a Meteor-MN2 soft-symbol file:
```
weatherdump lrpt ./file_path.s --jpeg
```

### Installation

### Future Features List
- [ ] Add multi-thread support to decoder (Beta 2).
- [ ] LRPT RGB component.
