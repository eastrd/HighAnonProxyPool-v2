from datetime import datetime


def log_print(message):
    '''Automatically adds timestamp before the message print
    Parameters:
    -----------
    message: String
        Message to be printed on the screen
    '''
    current_time = str(datetime.now().strftime("%m-%d %H:%M:%S"))
    print("%s \t %s" %(current_time, str(message)))