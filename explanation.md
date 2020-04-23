Para demonstrar o padrão factory method eu criei um gerenciador de conexões. Como motivação eu considerei uma aplicação de chat com diversas interfaces, que suporta clientes websocket, tcp socket, USB, Bluetooth, RS-232, etc... Para fazer esse chat, precisamos obviamente isolar a logica do chat dos diferentes protocolos suportados. Com essas condições, gostaríamos de ter um gerenciador de conexões, que sabe se conectar com qualquer dispositivo, que seja desacoplado da implementação dos diferentes tipos de conexão, permitindo facilidade de adicionar novos tipos de conexão. Ai que entra o factory method, o gerenciador de conexões defere a criação das conexões para outros módulos, permitindo que adicionemos mais tipos de conexão simplesmente adicionando módulos. 

Para isso, precisamos definir a interface entre esses módulos e o gerenciador de conexões. Então, o módulo do gerenciador de conexões vai conhecer duas interfaces, `ConnectionFactory` e `Connection`. Ele utiliza `ConnectionFactory` para criar conexões, que são expostas para os clientes do gerenciador. Portanto, `Connection` também é parte da interface do gerenciador para com seus clientes

- Falar da interface do connection manager

- Propor ideia de implementação do servidor de chat:
    - Imagina um modulo que implementa as conexões Bluetooth, ele vai implementar a interface `ConnectionFactory` com `BluetoothFactory`
    - Ele vai soltar uma thread que escuta por novos dispositivos bluetooth
    - Ao escutar uma nova conexão, ele cria uma instância de `BluetoothFactory`, com os parâmetros para conexão com o dispositivo
    - Ele adiciona essa fábrica ao gerenciador de conexões, linkando ela com o nome do dispositivo
    - Quando qualquer cliente do gerenciador de conexões quiser se conectar com esse dispositivo Bluetooth em particular, o gerenciador vai utilizar a instancia previamente registrada de `BluetoothFactory` será usada para criar uma conexão pro dispositivo