<script>
import { Radar } from 'vue-chartjs'
export default {
  extends: Radar,
  props: {
    matchData: {
      type: Array
    }
  },
  data () {
    return {
      chartOptions: {
        legend: {
          display: false
        },
        scale: {
          ticks: {
            display: false
          },
          pointLabels: {
            fontColor: '#ebdbb2'
          },
          gridLines: {
            color: 'rgba(235, 219, 178, 0.4)'
          },
          angleLines: {
            color: 'rgba(235, 219, 178, 0.4)'
          }
        },
      }
    }
  },
  mounted() {
    var adrs = {}
    var mapCounts = {}

    this.matchData.forEach(match => {
      if (!adrs[match.matchData.map]) {
        adrs[match.matchData.map] = 0
      }
      if (!mapCounts[match.matchData.map]) {
        mapCounts[match.matchData.map] = 0
      }

      adrs[match.matchData.map] += match.playerData.adr
      mapCounts[match.matchData.map] += 1
    })

    var adrData = {}

    Object.keys(adrs).forEach(key => {
      adrData[key] = adrs[key] / mapCounts[key]
    })

    this.renderChart({
      labels: Object.keys(adrData),
      datasets: [{
        label: 'ADR',
        backgroundColor: 'rgba(0, 119, 204, 0.2)',
        borderColor: 'rgba(0, 119, 204, 0.5)',
        data: Object.values(adrData)
      }]
    }, this.chartOptions)
  }
}
</script>