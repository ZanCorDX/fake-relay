# fake-relay
Extremely simple app to parse a block submission EXACTLY as flashbots relay does.

The proper way to do this would have been to refactor the relay a little to split data parsing from data processing allowing to just plug-in test code but this was too much work.

Based on https://github.com/flashbots/mev-boost-relay (at commit d8a0d7bdb5d135624cbd231fc041cadbd803ebd9) taking code from:
- mev-boost-relay/services/api/service.go handleSubmitNewBlock for the parsing
- mev-boost-relay/common/types_spec.go

Things you can do with this app:
- Check if you generated submissions are properly generated
- Record a  submissions to compare it latter with another one.
    Eg: Record a json submission and checks if the same submission using ssz matches.
    
- Ignore the app and continue with your life
