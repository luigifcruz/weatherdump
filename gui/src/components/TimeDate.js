import { CopyToClipboard } from 'react-copy-to-clipboard'
import { Clock as Analog } from 'react-clock'
import React, { Component } from 'react'
import * as rxa from '../redux/actions'
import { connect } from 'react-redux'
import { DateTime } from 'luxon'
import SunCalc from 'suncalc'
import '../styles/TimeDate.scss'

function getHour(date, timezone) {
    if (date.getTime() < 2000) {
        return "--:--:--";
    }
    return DateTime.fromJSDate(date).setZone(timezone).toFormat("HH:mm:ss");
}

function getDate(date, timezone) {
    return DateTime.fromJSDate(date).setZone(timezone).toFormat("DDDD");
}

function getISO(date, timezone) {
    return DateTime.fromJSDate(date).setZone(timezone).toISO();
}

class TimeDate extends Component {

    constructor(props) {
        super(props);
        this.state = {
            lat: 0.0,
            lon: 0.0,
            date: new Date(),
            suncalc: SunCalc.getTimes(0, 0, 0),
            clocks: [{
                timezone: "local",
                name: "Local Time"
            },{
                timezone: "utc",
                name: "Universal Time"
            },{
                timezone: "America/New_York",
                name: "New York Time"
            }]
        }
    }

    componentWillReceiveProps(newProps) {
        if (Math.abs(newProps.state.longitude - this.state.lon) > 0.0001 ||
            Math.abs(newProps.state.latitude - this.state.lat) > 0.0001) {
            this.setState({
                lat: this.props.state.latitude,
                lon: this.props.state.longitude, 
                suncalc: SunCalc.getTimes(
                    new Date(),
                    this.props.state.latitude,
                    this.props.state.longitude)
            });
        }
    }

    componentDidMount() {
        this.interval = setInterval(() => this.setState({ date: new Date() }), 1000);   
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    render() {
        return (
            <div className="TimeDate">
                <h2>Time & Date</h2>
                <div className="ClockSection">
                    <center>
                        {this.state.clocks.map(function(clock){
                            return (
                                <div className="ClockBlock">
                                    <h3>{clock.name}</h3>
                                    <Analog
                                        size={250}
                                        secondHandWidth={7}
                                        secondHandLength={80}
                                        minuteHandWidth={10}
                                        minuteMarksWidth={2}
                                        hourMarksWidth={4}
                                        hourHandWidth={10}
                                        value={getHour(this.state.date, clock.timezone)} />
                                    <CopyToClipboard text={getISO(this.state.date, clock.timezone)}>
                                        <div className="DigitalBlock">
                                            <h4>{getHour(this.state.date, clock.timezone)}</h4>
                                            <h5>{getDate(this.state.date, clock.timezone)}</h5>
                                        </div>
                                    </CopyToClipboard>
                                </div>
                            )
                        }, this)}
                    </center>
                </div>
                <h3>Solar Time</h3>
                    <div className="PredictionBox">
                        <p>Sunrise</p>
                        <h4>{getHour(this.state.suncalc.sunrise, "local")}</h4>
                        <label>When the top edge of the sun crosses the horizon.</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Sunrise End</p>
                        <h4>{getHour(this.state.suncalc.sunriseEnd, "local")}</h4>
                        <label>When the entire sun is above the horizon.</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Solar Noon</p>
                        <h4>{getHour(this.state.suncalc.solarNoon, "local")}</h4>
                        <label>When the sun reaches the higher elevation of the day.</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Golden Hour End</p>
                        <h4>{getHour(this.state.suncalc.goldenHour, "local")}</h4>
                        <label>When the best time for photography ends.</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Sunset Start</p>
                        <h4>{getHour(this.state.suncalc.sunsetStart, "local")}</h4>
                        <label>When the sun touches the horizon.</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Sunset</p>
                        <h4>{getHour(this.state.suncalc.sunset, "local")}</h4>
                        <label>When the sun disappears below the horizon.</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Dusk</p>
                        <h4>{getHour(this.state.suncalc.dusk, "local")}</h4>
                        <label>When the nautical twilight starts.</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Nautical Dusk</p>
                        <h4>{getHour(this.state.suncalc.nauticalDusk, "local")}</h4>
                        <label>When the evening astronomical twilight starts.</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Night Start</p>
                        <h4>{getHour(this.state.suncalc.night, "local")}</h4>
                        <label>When it's dark enough for astronomical observations.</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Nadir</p>
                        <h4>{getHour(this.state.suncalc.nadir, "local")}</h4>
                        <label>When the sun reaches the lowest position. Darkest time of the night!</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Night End</p>
                        <h4>{getHour(this.state.suncalc.nightEnd, "local")}</h4>
                        <label>When the morning astronomical twilight starts.</label>
                    </div>
                    <div className="PredictionBox">
                        <p>Dawn</p>
                        <h4>{getHour(this.state.suncalc.dawn, "local")}</h4>
                        <label>When the morning civil twilight starts.</label>
                    </div>
                <h3>Lunar Time</h3>
            </div>
        )
    }

}

TimeDate.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(TimeDate)  