const { getMatch } = require('./lib/StatmanClient')

async function main() {
    const response = await getMatch("111339799327333884")
    console.log(response)
}

main()