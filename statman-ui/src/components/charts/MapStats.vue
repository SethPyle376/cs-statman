<script>
import { Bar } from 'vue-chartjs'
export default {
  extends: Bar,
  props: {
    matchData: {
      type: Array
    }
  },
  data () {
    return {
      chartOptions: {
        legend: {
          labels: {
            fontColor: '#ebdbb2'
          }
        }
      }
    }
  },
  mounted() {
    var killMap = {}
    var deathMap = {}

    this.matchData.forEach(match => {
      console.log(match)
      if (!killMap[match.matchData.map]) {
        killMap[match.matchData.map] = 0
      }
      if (!deathMap[match.matchData.map]) {
        deathMap[match.matchData.map] = 0
      }
      killMap[match.matchData.map] += match.playerData.kills
      deathMap[match.matchData.map] += match.playerData.deaths
    })

    this.renderChart({
      labels: Object.keys(killMap),
      datasets: [{
        label: 'Kills',
        backgroundColor: 'rgba(152, 151, 26, 0.8)',
        data: Object.values(killMap)
      },{
        label: 'Deaths',
        backgroundColor: 'rgba(204, 36, 29, 0.8)',
        data: Object.values(deathMap)
      }]
    },
    this.chartOptions)
  }
}
</script>